# hash-file-compare
CLI tool that creates and compares file hashes, using SHA-256.

# Usage:

```python
hash-file-compare  (filename)
```
- returns hash value of (filename)


```python
hash-file-compare -d (directory)
```
- Scans through directory, displaying lists of duplicate files and their hash value

```python
hash-file-compare -TRASH (directory)
```
- (in progress) Scans through directory, moving all duplicate files to trash after the first found instance


```python
hash-file-compare -REMOVE (directory)
```
- Scans through directory, deleting all duplicate files after the first found instance
