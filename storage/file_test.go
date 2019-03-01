package storage

import "testing"

func Test_migrationFileNameRe(t *testing.T) {
	t.Run("bad migration file names", func(t *testing.T) {
		t.Parallel()
		badFileNames := map[string]string{
			"12 3-name-apply.sql": "spaces in versions are invalid",
			"123-n ame-appl.sql":  "spaces in names are invalid",
			"123-name.sql":        `"apply" or "rollback" is required`,
		}
		for name, reason := range badFileNames {
			if migratationFileNameRe.MatchString(name) {
				t.Errorf(`"%s" should be an invalid migration name because %s`, name, reason)
			}
		}
	})

	t.Run("good migration file names", func(t *testing.T) {
		t.Parallel()
		goodFileNames := [][4]string{
			{"001-name-apply.sql", "001", "name", "apply"},
			{"001_name_apply.sql", "001", "name", "apply"},
			{"001_name_is-here_apply.sql", "001", "name_is-here", "apply"},
			{"001-name-rollback.sql", "001", "name", "rollback"},
			{"001-n-a-m-e-rollback.sql", "001", "n-a-m-e", "rollback"},
			{"2021-12-01-10:13:13-name-rollback.sql", "2021-12-01-10:13:13", "name", "rollback"},
		}
		for _, goodCase := range goodFileNames {
			matches := migratationFileNameRe.FindStringSubmatch(goodCase[0])
			if matches[1] != goodCase[1] {
				t.Errorf("version should be %s for file name %s but got %s", goodCase[1], goodCase[1], matches[1])
			}
			if matches[2] != goodCase[2] {
				t.Errorf("migration name should be %s for file name %s but got %s", goodCase[2], goodCase[2], matches[2])
			}
			if matches[3] != goodCase[3] {
				t.Errorf("migration kind should be %s for file name %s but got %s", goodCase[3], goodCase[3], matches[3])
			}
		}
	})
}
