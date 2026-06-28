package m3u

import (
	"os"
)

// AtomicReplace replaces the old file with the new file. 
// If removeOriginal is true, it deletes the old file after successful replacement (though os.Rename does this).
// If removeOriginal is false, we might want to keep both? No, usually atomic replace means the new one replaces the old.
// The requirement says "-r / --remove-original equivalent". In many tools, this means "don't keep the source if it was a copy".
// But here we are talking about M3U files. 

// AtomicReplace moves the new file to the old path.
func AtomicReplace(newPath, oldPath string) error {
	return os.Rename(newPath, oldPath)
}

// In the context of "playlist hygiene", maybe it means if we are fixing extensions,
// do we want to keep the original .flac file or delete it? 
// The legacy script just changed the content of the M3U.

// Let's implement it as: if removeOriginal is true, we delete the source file.
func HandleRemoveOriginal(sourcePath string, remove bool) error {
	if remove {
		return os.Remove(sourcePath)
	}
	return nil
}
