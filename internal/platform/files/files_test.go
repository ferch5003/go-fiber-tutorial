package files

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFile_Successful(t *testing.T) {
	// Given
	file := "go.mod" // Go Mod file always exits.

	// When
	absoluteFilepath, err := GetFile(file)

	// Then
	require.NoError(t, err)
	require.FileExists(t, absoluteFilepath)
}

func TestGetFile_FailsDueToNotExistingFile(t *testing.T) {
	// Given
	file := "a_file_that_doesnt_exist"

	// When
	_, err := GetFile(file)

	// Then
	require.ErrorContains(t, err, "stat")
	require.ErrorContains(t, err, "a_file_that_doesnt_exist")
	require.ErrorContains(t, err, "no such file or directory")
}

func TestGetDir_Successful(t *testing.T) {
	// Given
	directory := "cmd" // CMD is the entrypoint of all applications.

	// When
	absoluteDirectory, err := GetDir(directory)

	// Then
	require.NoError(t, err)
	require.DirExists(t, absoluteDirectory)
}

func TestGetDir_FailsDueToNotExistingDirectory(t *testing.T) {
	// Given
	file := "a_directory_that_dont_exist"

	// When
	_, err := GetDir(file)

	// Then
	require.ErrorContains(t, err, "stat")
	require.ErrorContains(t, err, "a_directory_that_dont_exist")
	require.ErrorContains(t, err, "no such file or directory")
}
