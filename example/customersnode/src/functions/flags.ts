// Flags is a bitwise mask to determine is a flag is set in a frame.
export enum Flag {
  END_STREAM = 0x01,
  END_HEADERS = 0x04,
  PADDED = 0x08,
}

export function isFlagSet(f: Flag, flagsByte: Flag): boolean {
  return (flagsByte & f) !== 0;
}

export function hasFlag(f: Flag, other: Flag): boolean {
  return (other & f) !== 0;
}
