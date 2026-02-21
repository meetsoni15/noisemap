# noisemap demo tape
# Record with: vhs demo.tape

Output demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "TokyoNight"
Set Padding 20
Set TypingSpeed 80ms
Set PlaybackSpeed 0.8

# Launch noisemap on its own source
Type "noisemap ."
Sleep 500ms
Enter

# Wait for scan to complete
Sleep 6s

# Show file list — navigate up and down
Sleep 1s
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 1s

# Go to top, then bottom
Type "g"
Sleep 800ms
Type "G"
Sleep 800ms
Type "g"
Sleep 800ms

# Switch sort mode
Type "s"
Sleep 1s
Type "s"
Sleep 1s
Type "s"
Sleep 1s
Type "s"
Sleep 1s

# Switch to heatmap view
Type "v"
Sleep 2s

# Navigate in heatmap
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 400ms

# Back to list
Type "v"
Sleep 1.5s

# Switch pane to detail
Tab
Sleep 1.5s

# Go back to list pane
Tab
Sleep 1s

# Quit
Type "q"
Sleep 500ms
