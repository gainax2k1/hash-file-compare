# hashcomparefiles
CLI tool that computes and compares file hashes, using SHA-256. 
This tool also makes it easy to delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize. It uses SHA-256 to uniquely identify the file contents, so even if a duplicate file has a different name, it will still be flagged. The filesize is included for refrence, and for the remote chance of hash collision. 

* symlinks and empty files are ignored.

# Usage:

```python
hashcomparefiles (filename/directory)
```
- returns hash value of a single file or through directory and sub-directories, displaying lists of duplicate files with their size and their hash value

```python
hashcomparefiles -trash (directory)
```
- (Linux only) Scans through directory, moving all duplicate files to trash after the first found instance. Currently, only fully works on primary drive in Linux based systems. For non-Linux systems, a folder is created in the working directory, and files are move into that, with corresponding .trashinfo files being created to record original file location.
- (Currently, -trash uses os.Rename to trash the files, which might not work correctly on external mounts/devices. In these cases, it switches to a copy/delete process instead, which doesn not follow the FreeDesktop spec. I would like to improve this in the future.)
- In the future, will add ability to choose which file(s) to trash.


```python
hashcomparefiles -delete (directory)
```
- Scans through directory, deleting all duplicate files after the first found instance
- In the future, will add ability to choose which file(s) to delete.

```python
hashcomparefiles -p (directory)
```
- Scans through directory once without hashing to get total file count, then hashes the second run. Potentially useful for large runs (+1,000 files) to determine progress, at the cost of additional disk hits.

```python
hashcomparefiles -v (directory)
```
- Verbose - outputs every hash value and filesize, even if not duplicated.

```python
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory.

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions

```python
cat (filename) | hashcomparefiles -(flag)
```
- Pipe in list of files and or folders to compare against each other. Flags maintain functionality.

# Examples:

```python
hashcomparefiles testdata/

2026/03/06 18:48:58 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:48:58 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:48:58 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:48:58 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:48:58 Total files processed: 16


```


```python
hashcomparefiles -log default testdata/

2026/03/06 18:58:22 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:58:22 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:58:22 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:58:22 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:58:22 Total files processed: 16
gainax2k1@pop-os:~/Documents/workspace/hashcomparefiles$ more log.log 
2026/03/06 18:58:22 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:58:22 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:58:22 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:58:22 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:58:22 Total files processed: 16
```