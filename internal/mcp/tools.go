package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type EchoInput struct {
	Text string `json:"text" jsonschema:"text to echo"`
}

type EchoOutput struct {
	Text string `json:"text" jsonschema:"echoed text"`
}

type FSReadInput struct {
	Path string `json:"path" jsonschema:"relative file path"`
}

type FSReadOutput struct {
	Content string `json:"content" jsonschema:"file content"`
}

type FSWriteInput struct {
	Path    string `json:"path" jsonschema:"relative file path"`
	Content string `json:"content" jsonschema:"file content"`
}

type FSWriteOutput struct {
	OK bool `json:"ok" jsonschema:"write result"`
}

type FSListInput struct {
	Path      string `json:"path" jsonschema:"relative directory path"`
	Recursive bool   `json:"recursive" jsonschema:"recursive listing"`
}

type FSFileEntry struct {
	Path    string    `json:"path" jsonschema:"relative entry path"`
	IsDir   bool      `json:"is_dir" jsonschema:"is directory"`
	Size    int64     `json:"size" jsonschema:"entry size"`
	Mode    uint32    `json:"mode" jsonschema:"file mode"`
	ModTime time.Time `json:"mod_time" jsonschema:"modification time"`
}

type FSListOutput struct {
	Path    string        `json:"path" jsonschema:"listed path"`
	Entries []FSFileEntry `json:"entries" jsonschema:"entries"`
}

type FSStatInput struct {
	Path string `json:"path" jsonschema:"relative path"`
}

type FSStatOutput struct {
	Entry FSFileEntry `json:"entry" jsonschema:"entry"`
}

type FSDeleteInput struct {
	Path string `json:"path" jsonschema:"relative path"`
}

type FSDeleteOutput struct {
	OK bool `json:"ok" jsonschema:"delete result"`
}

type FSApplyPatchInput struct {
	Path  string `json:"path" jsonschema:"relative file path"`
	Patch string `json:"patch" jsonschema:"unified diff patch"`
}

type FSApplyPatchOutput struct {
	OK bool `json:"ok" jsonschema:"apply result"`
}

type FSMkdirInput struct {
	Path string `json:"path" jsonschema:"relative directory path"`
}

type FSMkdirOutput struct {
	OK bool `json:"ok" jsonschema:"mkdir result"`
}

type FSRenameInput struct {
	Source      string `json:"source" jsonschema:"relative source path"`
	Destination string `json:"destination" jsonschema:"relative destination path"`
}

type FSRenameOutput struct {
	OK bool `json:"ok" jsonschema:"rename result"`
}

type FSReadBase64Input struct {
	Path string `json:"path" jsonschema:"relative file path"`
}

type FSReadBase64Output struct {
	Data     string `json:"data" jsonschema:"base64-encoded file bytes"`
	MimeType string `json:"mime_type" jsonschema:"detected mime type"`
}

type GrepInput struct {
	Pattern string   `json:"pattern" jsonschema:"grep pattern"`
	Args    []string `json:"args" jsonschema:"grep options (flags only)"`
}

type GrepOutput struct {
	Stdout   string `json:"stdout" jsonschema:"grep standard output"`
	Stderr   string `json:"stderr" jsonschema:"grep standard error"`
	ExitCode int    `json:"exit_code" jsonschema:"grep exit code"`
}

func RegisterTools(server *sdkmcp.Server) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "echo", Description: "echo input text"}, echoTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.read", Description: "read file content"}, fsReadTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.read_base64", Description: "read file bytes as base64"}, fsReadBase64Tool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.write", Description: "write file content"}, fsWriteTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.list", Description: "list directory entries"}, fsListTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.stat", Description: "stat file or directory"}, fsStatTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.delete", Description: "delete file or directory"}, fsDeleteTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.apply_patch", Description: "apply unified diff patch"}, fsApplyPatchTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.mkdir", Description: "create directory (mkdir -p)"}, fsMkdirTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "fs.rename", Description: "rename/move file or directory"}, fsRenameTool)
	sdkmcp.AddTool(server, &sdkmcp.Tool{Name: "grep", Description: "grep within /data using GNU grep"}, grepTool)
}

func echoTool(ctx context.Context, req *sdkmcp.CallToolRequest, input EchoInput) (
	*sdkmcp.CallToolResult,
	EchoOutput,
	error,
) {
	return nil, EchoOutput{Text: input.Text}, nil
}

func fsReadTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSReadInput) (
	*sdkmcp.CallToolResult,
	FSReadOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSReadOutput{}, err
	}
	data, err := os.ReadFile(target)
	if err != nil {
		return nil, FSReadOutput{}, err
	}
	return nil, FSReadOutput{Content: string(data)}, nil
}

func fsReadBase64Tool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSReadBase64Input) (
	*sdkmcp.CallToolResult,
	FSReadBase64Output,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSReadBase64Output{}, err
	}
	data, err := os.ReadFile(target)
	if err != nil {
		return nil, FSReadBase64Output{}, err
	}
	ext := strings.ToLower(filepath.Ext(target))
	mimeType := mime.TypeByExtension(ext)
	if strings.TrimSpace(mimeType) == "" {
		// Fallback mapping for common image/audio extensions.
		switch ext {
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		case ".bmp":
			mimeType = "image/bmp"
		case ".svg":
			mimeType = "image/svg+xml"
		case ".mp3":
			mimeType = "audio/mpeg"
		case ".wav":
			mimeType = "audio/wav"
		case ".ogg":
			mimeType = "audio/ogg"
		case ".flac":
			mimeType = "audio/flac"
		default:
			mimeType = "application/octet-stream"
		}
	}
	return nil, FSReadBase64Output{
		Data:     base64.StdEncoding.EncodeToString(data),
		MimeType: mimeType,
	}, nil
}

func grepTool(ctx context.Context, req *sdkmcp.CallToolRequest, input GrepInput) (
	*sdkmcp.CallToolResult,
	GrepOutput,
	error,
) {
	if strings.TrimSpace(input.Pattern) == "" {
		return nil, GrepOutput{}, fmt.Errorf("pattern is required")
	}
	if stat, err := os.Stat("/data"); err != nil || !stat.IsDir() {
		return nil, GrepOutput{}, fmt.Errorf("/data is not available")
	}

	args := append([]string{}, input.Args...)
	args = append(args, input.Pattern, ".")

	cmd := exec.CommandContext(ctx, "grep", args...)
	cmd.Dir = "/data"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			if exitCode != 1 {
				return nil, GrepOutput{}, fmt.Errorf("grep failed: %s", strings.TrimSpace(stderr.String()))
			}
		} else {
			return nil, GrepOutput{}, err
		}
	}

	return nil, GrepOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

func fsWriteTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSWriteInput) (
	*sdkmcp.CallToolResult,
	FSWriteOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSWriteOutput{}, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return nil, FSWriteOutput{}, err
	}
	if err := os.WriteFile(target, []byte(input.Content), 0o644); err != nil {
		return nil, FSWriteOutput{}, err
	}
	return nil, FSWriteOutput{OK: true}, nil
}

func fsListTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSListInput) (
	*sdkmcp.CallToolResult,
	FSListOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePathAllowRoot(root, input.Path)
	if err != nil {
		return nil, FSListOutput{}, err
	}
	info, err := os.Stat(target)
	if err != nil {
		return nil, FSListOutput{}, err
	}
	if !info.IsDir() {
		return nil, FSListOutput{}, fmt.Errorf("path is not a directory")
	}

	entries := []FSFileEntry{}
	if input.Recursive {
		err = filepath.WalkDir(target, func(p string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if p == target {
				return nil
			}
			entryInfo, err := d.Info()
			if err != nil {
				return err
			}
			entry, err := entryForPath(root, p, entryInfo)
			if err != nil {
				return err
			}
			entries = append(entries, entry)
			return nil
		})
	} else {
		dirEntries, err := os.ReadDir(target)
		if err != nil {
			return nil, FSListOutput{}, err
		}
		for _, entry := range dirEntries {
			entryInfo, err := entry.Info()
			if err != nil {
				return nil, FSListOutput{}, err
			}
			fullPath := filepath.Join(target, entry.Name())
			fileEntry, err := entryForPath(root, fullPath, entryInfo)
			if err != nil {
				return nil, FSListOutput{}, err
			}
			entries = append(entries, fileEntry)
		}
	}
	if err != nil {
		return nil, FSListOutput{}, err
	}

	listedPath := strings.TrimSpace(input.Path)
	if listedPath == "" {
		listedPath = "."
	}
	return nil, FSListOutput{Path: listedPath, Entries: entries}, nil
}

func fsStatTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSStatInput) (
	*sdkmcp.CallToolResult,
	FSStatOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePathAllowRoot(root, input.Path)
	if err != nil {
		return nil, FSStatOutput{}, err
	}
	info, err := os.Stat(target)
	if err != nil {
		return nil, FSStatOutput{}, err
	}
	entry, err := entryForPath(root, target, info)
	if err != nil {
		return nil, FSStatOutput{}, err
	}
	return nil, FSStatOutput{Entry: entry}, nil
}

func fsDeleteTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSDeleteInput) (
	*sdkmcp.CallToolResult,
	FSDeleteOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSDeleteOutput{}, err
	}
	if err := os.RemoveAll(target); err != nil {
		return nil, FSDeleteOutput{}, err
	}
	return nil, FSDeleteOutput{OK: true}, nil
}

func fsApplyPatchTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSApplyPatchInput) (
	*sdkmcp.CallToolResult,
	FSApplyPatchOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSApplyPatchOutput{}, err
	}
	orig, err := os.ReadFile(target)
	if err != nil {
		return nil, FSApplyPatchOutput{}, err
	}
	updated, err := applyUnifiedPatch(string(orig), input.Patch)
	if err != nil {
		return nil, FSApplyPatchOutput{}, err
	}
	info, err := os.Stat(target)
	if err != nil {
		return nil, FSApplyPatchOutput{}, err
	}
	if err := os.WriteFile(target, []byte(updated), info.Mode().Perm()); err != nil {
		return nil, FSApplyPatchOutput{}, err
	}
	return nil, FSApplyPatchOutput{OK: true}, nil
}

func fsMkdirTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSMkdirInput) (
	*sdkmcp.CallToolResult,
	FSMkdirOutput,
	error,
) {
	root := dataRoot()
	target, err := resolvePath(root, input.Path)
	if err != nil {
		return nil, FSMkdirOutput{}, err
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return nil, FSMkdirOutput{}, err
	}
	return nil, FSMkdirOutput{OK: true}, nil
}

func fsRenameTool(ctx context.Context, req *sdkmcp.CallToolRequest, input FSRenameInput) (
	*sdkmcp.CallToolResult,
	FSRenameOutput,
	error,
) {
	root := dataRoot()
	source, err := resolvePath(root, input.Source)
	if err != nil {
		return nil, FSRenameOutput{}, err
	}
	destination, err := resolvePath(root, input.Destination)
	if err != nil {
		return nil, FSRenameOutput{}, err
	}

	if _, err := os.Lstat(destination); err == nil {
		return nil, FSRenameOutput{}, fmt.Errorf("destination already exists")
	} else if !os.IsNotExist(err) {
		return nil, FSRenameOutput{}, err
	}

	if err := os.Rename(source, destination); err != nil {
		return nil, FSRenameOutput{}, err
	}
	return nil, FSRenameOutput{OK: true}, nil
}

func dataRoot() string {
	root := strings.TrimSpace(os.Getenv("MCP_DATA_DIR"))
	if root == "" {
		root = "/data"
	}
	return root
}

func resolvePathAllowRoot(root, requestPath string) (string, error) {
	if strings.TrimSpace(requestPath) == "" {
		return root, nil
	}
	return resolvePath(root, requestPath)
}

func resolvePath(root, requestPath string) (string, error) {
	clean := filepath.Clean(requestPath)
	if clean == "." || clean == "" {
		return "", os.ErrInvalid
	}
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", os.ErrInvalid
	}
	return filepath.Join(root, clean), nil
}

func entryForPath(root, target string, info os.FileInfo) (FSFileEntry, error) {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return FSFileEntry{}, err
	}
	if strings.HasPrefix(rel, "..") {
		return FSFileEntry{}, os.ErrInvalid
	}
	if rel == "." {
		rel = ""
	}
	return FSFileEntry{
		Path:    filepath.ToSlash(rel),
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		Mode:    uint32(info.Mode().Perm()),
		ModTime: info.ModTime(),
	}, nil
}

func applyUnifiedPatch(original, patch string) (string, error) {
	lines := strings.Split(original, "\n")
	out := make([]string, 0, len(lines))
	index := 0
	patchLines := strings.Split(patch, "\n")
	hunksApplied := 0

	for i := 0; i < len(patchLines); i++ {
		line := patchLines[i]
		if !strings.HasPrefix(line, "@@") {
			continue
		}

		origStart, err := parseUnifiedHunkHeader(line)
		if err != nil {
			return "", err
		}
		origStart--
		if origStart < 0 {
			origStart = 0
		}
		if origStart > len(lines) {
			return "", fmt.Errorf("patch out of range")
		}

		out = append(out, lines[index:origStart]...)
		index = origStart
		hunksApplied++

		for i+1 < len(patchLines) {
			next := patchLines[i+1]
			if strings.HasPrefix(next, "@@") {
				break
			}
			i++

			if next == "" {
				if i == len(patchLines)-1 {
					break
				}
				return "", fmt.Errorf("invalid patch line")
			}
			if next[0] == '\\' {
				continue
			}
			op := next[0]
			text := next[1:]
			switch op {
			case ' ':
				if index >= len(lines) || lines[index] != text {
					return "", fmt.Errorf("patch context mismatch")
				}
				out = append(out, text)
				index++
			case '-':
				if index >= len(lines) || lines[index] != text {
					return "", fmt.Errorf("patch delete mismatch")
				}
				index++
			case '+':
				out = append(out, text)
			default:
				return "", fmt.Errorf("invalid patch operation")
			}
		}
	}
	if hunksApplied == 0 {
		return "", fmt.Errorf("patch contains no hunks")
	}

	out = append(out, lines[index:]...)
	return strings.Join(out, "\n"), nil
}

func parseUnifiedHunkHeader(header string) (int, error) {
	trimmed := strings.TrimPrefix(header, "@@")
	trimmed = strings.TrimSpace(trimmed)
	if !strings.HasPrefix(trimmed, "-") {
		return 0, fmt.Errorf("invalid hunk header")
	}
	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid hunk header")
	}

	origPart := strings.TrimPrefix(parts[0], "-")
	origFields := strings.SplitN(origPart, ",", 2)
	origStart, err := strconv.Atoi(origFields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hunk header")
	}
	return origStart, nil
}
