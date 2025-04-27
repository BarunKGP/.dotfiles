#!/bin/bash

SAVE_FILE="$HOME/.dotfiles/tmux/.tmux_sessions"

save_sessions() {
    echo "Saving tmux sessions..."
    rm -f "$SAVE_FILE"

    tmux list-sessions -F "#{session_name}" | while read -r session; do
        echo "Session: $session" >>"$SAVE_FILE"
        tmux list-windows -t "$session" -F "#{window_index}:#{window_name}" | while read -r window; do
            echo "  Window: $window" >>"$SAVE_FILE"
            tmux list-panes -t "$session" -F "#{window_index}:#{pane_index}:#{pane_current_path}" | while read -r pane; do
                echo "    Pane: $pane" >>"$SAVE_FILE"
            done
        done
    done

    echo "Sessions saved to $SAVE_FILE"
}

restore_sessions() {
    if [ ! -f "$SAVE_FILE" ]; then
        echo "No saved sessions found."
        return 1
    fi

    echo "Restoring tmux sessions..."
    while IFS= read -r line; do
        case "$line" in
        "Session: "*)
            session_name="${line#Session: }"
            tmux new-session -d -s "$session_name"
            ;;
        "  Window: "*)
            window_info="${line#  Window: }"
            window_index="${window_info%%:*}"
            window_name="${window_info#*:}"
            tmux new-window -t "$session_name:$window_index" -n "$window_name"
            ;;
        "    Pane: "*)
            pane_info="${line#    Pane: }"
            window_index=$(echo "$pane_info" | cut -d':' -f1)
            pane_index=$(echo "$pane_info" | cut -d':' -f2)
            pane_path=$(echo "$pane_info" | cut -d':' -f3-)

            # Navigate to the correct working directory
            if [ "$pane_index" -eq 0 ]; then
                tmux send-keys -t "$session_name:$window_index" "cd '$pane_path'" C-m
            else
                tmux split-window -t "$session_name:$window_index" -c "$pane_path"
            fi
            ;;
        esac
    done <"$SAVE_FILE"

    echo "Sessions restored."
}

case "$1" in
save) save_sessions ;;
restore) restore_sessions ;;
*) echo "Usage: $0 {save|restore}" ;;
esac
