# voxel

Voxel-fun 3D visualizer integration.

https://github.com/pawel-dubiel/voxel-fun

## Why

Provides an alternative 3D visualization for the drone simulator using voxel graphics.

## How

```sh
task voxel:start        # Start voxel-fun visualizer
task voxel:stop         # Stop visualizer
task voxel:deps:install # Clone and install dependencies
task voxel:deps:clean   # Remove source directory
task voxel:debug        # Print debug info
```

## Configuration

- `VOXEL_PORT` - Dev server port (default: 5173)

Source is cloned to `.src/voxel-fun/`.
