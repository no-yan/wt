# wt: リファクタリング計画

## 概要

コード品質レビューで発見された設計上の問題を修正する。
「A Philosophy of Software Design」の原則（情報隠蔽、Deep Module、一貫性）に基づく。

## 問題と解決策

### 1. VCS詳細のコマンド層への漏洩 [Critical]

**現状**: `cmd/add.go:42-44` でVCSタイプに応じた条件分岐

```go
opts := internal.AddOptions{Revision: addRevision}
if vcsType == internal.VCSGit {
    opts.Branch = name // For Git, name is the branch
}
```

**解決策**: `AddOptions` を簡素化し、VCS固有ロジックをManager内部に移動

```go
// Before
type AddOptions struct {
    Branch   string  // Git用
    Revision string  // jj用
}

// After
type AddOptions struct {
    Revision string  // 両VCS共通: ベースとなるリビジョン（省略可）
}
```

**変更ファイル**:
- `internal/vcs.go`: `AddOptions` からBranchを削除
- `internal/worktree_manager.go`: `AddWorkspace` でnameをbranch名として使用
- `cmd/add.go`: VCSタイプ判定を削除

---

### 2. Worktree/Workspace 型の二重管理 [Critical]

**現状**: Git専用の `Worktree` 型と統一された `Workspace` 型が混在

| ファイル | 問題 |
|---------|------|
| `internal/worktree.go` | `Worktree` 型定義 |
| `internal/git_service.go` | `ListWorktrees()` と `ListWorkspaces()` 両方存在 |
| `cmd/list.go` | レガシー関数 `filterWorktrees`, `formatWorktreeList` が残存 |

**解決策**: `Worktree` 型を廃止し `Workspace` に統一

**変更ファイル**:
- `internal/worktree.go`: `Worktree` 型を削除、`Status` のみ残す
- `internal/git_service.go`: `ListWorktrees()` を削除
- `internal/git.go`: `ParseWorktreeList` を `[]Workspace` を返すよう変更
- `internal/worktree_manager.go`: `Workspace` 型を使用
- `cmd/list.go`: レガシー関数を削除

---

### 3. パススルーメソッドの排除 [Important]

**現状**: `internal/jj_service.go:115-117`

```go
func (j *JjService) Run(command string) (string, error) {
    return j.runner.Run(command)
}
```

**解決策**: `JjWorkspaceManager` が `CommandRunner` を直接保持

```go
type JjWorkspaceManager struct {
    jjService *JjService
    runner    CommandRunner  // 追加
}
```

**変更ファイル**:
- `internal/jj_service.go`: `Run()` メソッドを削除
- `internal/jj_workspace_manager.go`: runnerを直接保持
- `internal/vcs_detect.go`: ファクトリでrunnerを渡す

---

### 4. バリデーション関数の統一 [Important]

**現状**: ほぼ同一のバリデーション関数が2つ存在

```go
// jj_workspace_manager.go
func validateWorkspaceName(name string) error { ... }

// worktree_manager.go
func validateBranchName(branch string) error { ... }
```

**解決策**: 共通関数に統一

```go
// internal/validation.go (新規)
func ValidateName(name string) error {
    if name == "" {
        return fmt.Errorf("name cannot be empty")
    }
    if strings.ContainsAny(name, ";&|`$(){}[]<>") {
        return fmt.Errorf("invalid characters in name: %q", name)
    }
    if strings.HasPrefix(name, "-") {
        return fmt.Errorf("name cannot start with dash: %q", name)
    }
    return nil
}

func ValidatePath(path string) error { ... }
```

**変更ファイル**:
- `internal/validation.go`: 新規作成
- `internal/worktree_manager.go`: `validateBranchName`, `validatePath` を削除
- `internal/jj_workspace_manager.go`: `validateWorkspaceName` を削除

---

## 実装順序

```
1. バリデーション関数統一 (独立して実施可能)
   └── internal/validation.go 作成
   └── 既存関数を置換

2. Worktree/Workspace 型統一
   └── internal/worktree.go 修正
   └── internal/git.go 修正
   └── internal/git_service.go 修正
   └── internal/worktree_manager.go 修正
   └── cmd/list.go レガシー関数削除

3. パススルーメソッド排除
   └── internal/jj_workspace_manager.go 修正
   └── internal/jj_service.go 修正
   └── internal/vcs_detect.go 修正

4. AddOptions 簡素化
   └── internal/vcs.go 修正
   └── internal/worktree_manager.go 修正
   └── cmd/add.go 修正
```

## 期待される成果

- コマンド層がVCS詳細を知らなくて済む
- 型が統一され認知負荷が下がる
- 重複コードが削減される
- 将来のVCS追加（Mercurial等）が容易になる
