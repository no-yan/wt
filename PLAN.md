# wt: Dual VCS Support (Git + jj) 実装計画

## 概要

wtにJujutsu (jj) workspace サポートを追加し、Git worktreeと両立させる。

**設計方針**:
- 過度な抽象化を避け、各VCSのネイティブワークフローを尊重
- Git: 既存の動作を維持 (1 worktree = 1 branch)
- jj: シンプルなworkspace管理 (bookmark強制なし)

## アーキテクチャ

```
cmd/                         # CLI層（VCS検出 → 適切なサービス選択）
    ↓
internal/vcs.go             # VCSService/WorkspaceManager インターフェース
    ↓
internal/git_service.go     # Git実装（既存を拡張）
internal/jj_service.go      # jj実装（新規）
```

## 新規ファイル

### 1. `internal/vcs.go` - 共通インターフェース

```go
type VCSType string
const (
    VCSGit VCSType = "git"
    VCSJj  VCSType = "jj"
)

// Workspace - 統一表現（Git worktree/jj workspace共通）
type Workspace struct {
    Name      string
    Path      string
    Revision  string   // Git: HEAD commit, jj: change ID
    Status    Status
    Branch    string   // Git only: branch名（jjでは空）
    Bookmarks []string // jj only: trunk()..@範囲のブックマーク名
}

type VCSService interface {
    Type() VCSType
    ListWorkspaces() ([]Workspace, error)
    GetDetailedStatus(path string) ([]string, error)
}

type WorkspaceManager interface {
    AddWorkspace(repoPath, name string, opts AddOptions) (string, error)
    RemoveWorkspace(repoPath, name string) error
}

type AddOptions struct {
    Branch   string  // Git用
    Revision string  // jj用
}
```

### 2. `internal/vcs_detect.go` - VCS自動検出

```go
func DetectVCS(path string) (VCSType, error)
func GetRepoRoot(runner CommandRunner, vcsType VCSType) (string, error)
func CheckVCSAvailable(vcsType VCSType) error  // jj/git がインストール済みか確認
```

- `.jj/` を優先（colocatedリポジトリ対応）
- 親ディレクトリを遡って検索
- **エラーハンドリング**: jj未インストール時は明確なエラーメッセージで fail fast

### 3. `internal/jj_service.go` - jjサービス実装

```go
type JjService struct { runner CommandRunner }

func (j *JjService) ListWorkspaces() ([]Workspace, error)
// jj workspace list をパース
// 各workspaceについて、trunk()..@ の範囲にあるブックマークを取得:
//   jj log -r 'trunk()..@' --no-graph -T 'bookmarks'
// ブックマークがあれば Workspace.Bookmarks に設定

func (j *JjService) GetDetailedStatus(path string) ([]string, error)
// jj status をパース

func (j *JjService) GetBookmarksInRange(workspacePath, range string) ([]string, error)
// 指定範囲のブックマークを取得
```

### 4. `internal/jj_workspace_manager.go` - jj workspace操作

```go
type JjWorkspaceManager struct { ... }

func (m *JjWorkspaceManager) AddWorkspace(repoPath, name string, opts AddOptions) (string, error)
// jj workspace add <path> --name <name> [--revision <rev>]

func (m *JjWorkspaceManager) RemoveWorkspace(repoPath, name string) error
// jj workspace forget + ディレクトリ削除
```

### 5. `internal/jj.go` - jj出力パース

```go
func ParseJjWorkspaceList(output string) []Workspace
func ParseJjStatus(output string) Status
```

## 修正ファイル

### `internal/git_service.go`
- `VCSService` インターフェース実装を追加
- `Type() VCSType` メソッド追加
- `ListWorktrees()` → `ListWorkspaces()` としても使えるように

### `internal/worktree_manager.go`
- `WorkspaceManager` インターフェース実装を追加

### `cmd/add.go`
```go
// 変更前: wt add <branch>
// 変更後: wt add <name> [-r/--revision <rev>]

// VCS検出 → 適切なmanager選択
vcsType := internal.DetectVCS(".")
switch vcsType {
case internal.VCSGit:
    // 既存動作: name=branch
case internal.VCSJj:
    // name=workspace名, -r=base revision
}
```

### `cmd/list.go`
- `VCSService.ListWorkspaces()` を使用
- jjの場合:
  - Change ID を表示
  - `trunk()..@` の範囲にブックマークがあれば、ブックマーク名も表示
  - 例: `abc123 (feature-auth)` または `abc123` (ブックマークなし)

### `cmd/switch.go`
- 変更なし（workspace名→パス解決は共通）

### `cmd/remove.go`
- `WorkspaceManager.RemoveWorkspace()` を使用

### `cmd/clean.go`
- Git: `git worktree prune`
- jj: stale workspace検出 + 削除

### `cmd/shell_init.go`
- jj用の補完追加

## コマンド動作の違い

| コマンド | Git | jj |
|---------|-----|-----|
| `wt add foo` | branch `foo` + worktree作成 | workspace `foo` 作成 |
| `wt add foo -r main` | (無視) | `main` ベースでworkspace作成 |
| `wt list` | Branch列表示 | Change ID + ブックマーク名表示 |
| `wt remove foo` | `git worktree remove` | `jj workspace forget` + rm |

### `wt list` 出力例 (jj)

```
NAME        CHANGE    BOOKMARK       STATUS
default     abc12345  (main)         clean
my-feature  def67890  feature-auth   dirty
experiment  ghi11111                 clean
```

- CHANGE: 短縮change ID
- BOOKMARK: `trunk()..@` 範囲にあるブックマーク（なければ空）

## 実装順序

1. **Phase 1: インターフェース定義**
   - `internal/vcs.go` 作成
   - `internal/vcs_detect.go` 作成

2. **Phase 2: Git側をインターフェースに適合**
   - `internal/git_service.go` 修正
   - `internal/worktree_manager.go` 修正

3. **Phase 3: jj実装**
   - `internal/jj.go` 作成
   - `internal/jj_service.go` 作成
   - `internal/jj_workspace_manager.go` 作成

4. **Phase 4: CLI層更新**
   - `cmd/add.go` 修正 (`-r` フラグ追加)
   - `cmd/list.go` 修正
   - `cmd/remove.go` 修正
   - `cmd/clean.go` 修正

5. **Phase 5: Shell integration**
   - `cmd/shell_init.go` 修正

6. **Phase 6: テスト**
   - 各新規ファイルのユニットテスト（MockCommandRunnerを使用）
   - jjインテグレーションテスト（jjがインストールされている環境のみ）
   - CI: jj未インストール環境ではjj関連テストをスキップ（`+build jj` タグ or 環境検出）

## 重要ファイルパス

**新規作成:**
- `internal/vcs.go`
- `internal/vcs_detect.go`
- `internal/jj.go`
- `internal/jj_service.go`
- `internal/jj_workspace_manager.go`

**修正:**
- `internal/git_service.go`
- `internal/worktree_manager.go`
- `cmd/add.go`
- `cmd/list.go`
- `cmd/remove.go`
- `cmd/clean.go`
- `cmd/shell_init.go`
