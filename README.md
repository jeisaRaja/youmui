# youmui

**youmui** is a command-line interface (TUI) application written in Go that allows users to search, play, and manage YouTube music playlists. It provides a sleek, user-friendly interface for browsing songs, adding them to a queue, and controlling playback—all from the comfort of your terminal.

## Features

- **Search for Songs**: Quickly find songs on YouTube using a search bar.
- **Song Management**: Add songs to a queue and play them back.
- **Playlist Support**: Manage and access your playlists easily.
- **Responsive UI**: A smooth and responsive terminal user interface built with the Charmbracelet Bubble Tea framework.
- **Audio Playback**: Stream audio directly from YouTube with seamless playback controls.

## Installation

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)
- `yt-dlp` for audio streaming. You can install it using:

  ```bash
  pip install -U yt-dlp
  ```

### Clone the Repository

```bash
git clone https://github.com/jeisaraja/youmui.git
cd youmui
```

### Build the Application

```bash
go build
```

### Run the Application

```bash
./youmui
```

## Usage

1. **Navigating Tabs**:

   - Press `s` to switch to the **Song** tab.
   - Press `p` to switch to the **Playlist** tab.
   - Press `q` to switch to the **Queue** tab.

2. **Searching for Songs**:

   - Press `f` to enter the search mode.
   - Type your search query and press `Enter`.

3. **Controlling Playback**:

   - Press `x` to play or pause the currently playing song.
   - Press `n` to play the next song in the queue.
   - Use `-` and `=` to adjust the volume.

4. **Managing the Queue**:

   - Select a song from the list and press `a` to add it to the queue.

5. **Exiting the Application**:
   - Press `Ctrl+C` to exit the application at any time.

## Contributing

Contributions are welcome! If you have suggestions for improvements or want to report a bug, feel free to open an issue or submit a pull request.
