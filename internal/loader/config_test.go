package loader

import "testing"

func TestMergeStackConfig(t *testing.T) {
	global := &Config{
		CommonArgs: []string{"--compatibility"},
		UpArgs:     []string{"--detach"},
		DownArgs:   []string{"--remove-orphans"},
	}

	t.Run("empty stack config uses global", func(t *testing.T) {
		merged := global.MergeStackConfig(&StackConfig{})
		assertSliceEqual(t, merged.CommonArgs, global.CommonArgs)
		assertSliceEqual(t, merged.UpArgs, global.UpArgs)
		assertSliceEqual(t, merged.DownArgs, global.DownArgs)
	})

	t.Run("stack config overrides up-args", func(t *testing.T) {
		stack := &StackConfig{UpArgs: []string{"--no-build"}}
		merged := global.MergeStackConfig(stack)
		assertSliceEqual(t, merged.UpArgs, []string{"--no-build"})
		assertSliceEqual(t, merged.DownArgs, global.DownArgs) // unchanged
	})

	t.Run("stack config overrides down-args", func(t *testing.T) {
		stack := &StackConfig{DownArgs: []string{"--volumes"}}
		merged := global.MergeStackConfig(stack)
		assertSliceEqual(t, merged.DownArgs, []string{"--volumes"})
		assertSliceEqual(t, merged.UpArgs, global.UpArgs) // unchanged
	})

	t.Run("original config not modified", func(t *testing.T) {
		original := len(global.UpArgs)
		_ = global.MergeStackConfig(&StackConfig{UpArgs: []string{"--changed"}})
		if len(global.UpArgs) != original {
			t.Error("Global config was modified")
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("defaults without config file", func(t *testing.T) {
		config, err := LoadConfig(t.TempDir())
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if len(config.UpArgs) != 5 {
			t.Errorf("Expected 5 default up-args, got %d", len(config.UpArgs))
		}
		if config.Timeout != 10 {
			t.Errorf("Expected default timeout=10, got %d", config.Timeout)
		}
	})

	t.Run("loads from yaml file", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "config.yaml", `up-args: ["--detach"]`)

		config, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		assertSliceEqual(t, config.UpArgs, []string{"--detach"})
	})

	t.Run("adds env-file when .env exists", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, ".env", "KEY=value")

		config, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if len(config.CommonArgs) != 2 || config.CommonArgs[0] != "--env-file" {
			t.Errorf("Expected --env-file in common-args, got %v", config.CommonArgs)
		}
	})

	t.Run("invalid yaml returns error", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "config.yaml", "invalid: yaml: [")

		_, err := LoadConfig(dir)
		if err == nil {
			t.Error("Expected error for invalid YAML")
		}
	})
}

func TestLoadStackConfig(t *testing.T) {
	t.Run("empty without config file", func(t *testing.T) {
		config, err := LoadStackConfig(t.TempDir())
		if err != nil {
			t.Fatalf("LoadStackConfig failed: %v", err)
		}
		if len(config.UpArgs) != 0 || len(config.DownArgs) != 0 {
			t.Error("Expected empty config")
		}
	})

	t.Run("loads from yaml file", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "config.yaml", `up-args: ["--no-build"]`)

		config, err := LoadStackConfig(dir)
		if err != nil {
			t.Fatalf("LoadStackConfig failed: %v", err)
		}
		assertSliceEqual(t, config.UpArgs, []string{"--no-build"})
	})
}
