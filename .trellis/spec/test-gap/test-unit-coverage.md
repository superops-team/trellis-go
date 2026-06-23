# test-unit-coverage: 单元测试覆盖率提升

> 针对 trellis-go 覆盖率短板包，补充单元测试。

## 当前状态

| 包 | 覆盖率 | 问题 |
|----|:---:|------|
| pkg/task | 44.9% | 大量边界路径未覆盖 |
| pkg/git | 36.0% | Git 操作 mock 不足 |
| pkg/prd | 0.0% | 完全没有测试 |
| pkg/fsutil | 77.0% | 部分函数未覆盖 |
| pkg/template | 76.4% | 错误路径未覆盖 |
| pkg/skill | 72.5% | 边界情况未覆盖 |
| pkg/upgrade | 72.3% | 网络请求 mock 不足 |

## 需求

### 1. pkg/task 覆盖率 → 80%+

在 `pkg/task/manager_test.go` 中补充：

```go
func TestManager_CreateTaskWithEmptyName(t *testing.T) { ... }
func TestManager_CreateTaskWithSpecialChars(t *testing.T) { ... }
func TestManager_StartTaskWithoutPRD(t *testing.T) { ... }
func TestManager_ArchiveTaskNotStarted(t *testing.T) { ... }
func TestManager_ArchiveTaskAlreadyArchived(t *testing.T) { ... }
func TestManager_GetTaskNotFound(t *testing.T) { ... }
func TestManager_ListTasksEmpty(t *testing.T) { ... }
func TestManager_ListTasksWithFilter(t *testing.T) { ... }
func TestManager_CurrentTaskNone(t *testing.T) { ... }
func TestManager_CurrentTaskMultiple(t *testing.T) { ... }
```

### 2. pkg/git 覆盖率 → 60%+

在 `pkg/git/git_test.go` 中补充：

```go
func TestGit_IsRepoTrue(t *testing.T) { ... }
func TestGit_IsRepoFalse(t *testing.T) { ... }
func TestGit_CurrentBranch(t *testing.T) { ... }
func TestGit_LatestCommit(t *testing.T) { ... }
func TestGit_HasUncommittedChanges(t *testing.T) { ... }
func TestGit_RootDir(t *testing.T) { ... }
```

### 3. pkg/prd 覆盖率 → 80%+

新建 `pkg/prd/prd_test.go`：

```go
func TestPRD_LoadValid(t *testing.T) { ... }
func TestPRD_LoadEmpty(t *testing.T) { ... }
func TestPRD_LoadMissing(t *testing.T) { ... }
func TestPRD_ValidateRequiredSections(t *testing.T) { ... }
func TestPRD_ValidateMissingSection(t *testing.T) { ... }
func TestPRD_ParseMarkdown(t *testing.T) { ... }
```

### 4. pkg/fsutil 覆盖率 → 90%+

```go
func TestFsutil_EnsureDir(t *testing.T) { ... }
func TestFsutil_EnsureDirNested(t *testing.T) { ... }
func TestFsutil_CopyFile(t *testing.T) { ... }
func TestFsutil_CopyFileOverwrite(t *testing.T) { ... }
func TestFsutil_IsSymlink(t *testing.T) { ... }
```

### 5. pkg/template 覆盖率 → 85%+

```go
func TestEngine_RenderWithMissingVar(t *testing.T) { ... }
func TestEngine_RenderWithNilContext(t *testing.T) { ... }
func TestEngine_RenderConditionalFalse(t *testing.T) { ... }
func TestEngine_RenderLoop(t *testing.T) { ... }
```

### 6. pkg/skill 覆盖率 → 85%+

```go
func TestSkill_LoadInvalidFormat(t *testing.T) { ... }
func TestSkill_LoadMissingName(t *testing.T) { ... }
func TestSkill_LoadEmptyFile(t *testing.T) { ... }
func TestSkill_FormatValidation(t *testing.T) { ... }
```

### 7. pkg/upgrade 覆盖率 → 85%+

```go
func TestChecker_NetworkError(t *testing.T) { ... }
func TestChecker_InvalidResponse(t *testing.T) { ... }
func TestChecker_VersionComparison(t *testing.T) { ... }
```

## 验收标准

| 包 | 当前 | 目标 |
|----|:---:|:---:|
| pkg/task | 44.9% | ≥80% |
| pkg/git | 36.0% | ≥60% |
| pkg/prd | 0.0% | ≥80% |
| pkg/fsutil | 77.0% | ≥90% |
| pkg/template | 76.4% | ≥85% |
| pkg/skill | 72.5% | ≥85% |
| pkg/upgrade | 72.3% | ≥85% |

## 参考

- trellis-go 现有测试：`pkg/*/`
- 覆盖率报告：`go test -cover ./pkg/...`
