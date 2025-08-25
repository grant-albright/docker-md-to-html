# docker-md-to-html
This is a container that'll watch a folder for markdown files and then convert them to html into a different folder.

## Metadata
Metadata is important to markdown. Wrap your markdown's metadata with 5 dashes at the beginning of the file. When converting to html, the metadata section will be ignored.

```md
-----
Title: Example of Using Metadata
Description: Insert description here.
-----
# Example of Using Metadata
This has been an example using metadata!
```
