# Hexdunk
A simple gui hex-editor with vi bindings

## Features
- Fast loading and editing of arbitrary large files (currently limited to about
  5GB because of problems with the imgui scrolling implementation. Will be
  fixed later to use home-made scrolling)
- Fast cut/copy/paste (in-program; not using system clipboard)
- Looks good (uses the my dear imgui library)
- Fully written in Go

It uses the excellent packages:
- "filebuf" (https://github.com/snhmibby/filebuf) for fast file loading and editing
- "filetree" (https://github.com/snhmibby/filetree) for a beautiful file-system dialog.

## Manual
The following operations are supported in the hex-window:
- left click: select byte with cursor
- right click: edit menu popup
- mouse dragging: select bytes
- shift clicking: extends selection

The following keys are bound (vi like)
- h,j,k,l: move around
- arrow keys: move around
- i: insert mode (insert bytes before the cursor)
- o: overwrite mode (overwrite bytes)
- escape: normal mode (move around/editing operations)
- x: cut
- y: copy
- p: paste

## Installation
go get github.com/snhmibby/hexdunk@main

## Warning
This software is (very much) in alpha version-state. It works for simple use cases,
but don't try to do complicated things.
I.e. there is a bug where copying from an edited file, pasting in another,
then writing the edited (1st) file to disk, could potentially change the other file-view. So keep all your write-back-to disk edits to 1 file for now :)

## Screenshots

![Image of HexDunk editing a selection](screenshots/selection_with_edit_menu.png)
![Image of the file dialog (proud of my work :X)](screenshots/open-dialog.png)
