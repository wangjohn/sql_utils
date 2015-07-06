package schema

func ReadSQL(filename string) ([]string, error) {
  var stmts []string

  f, err := os.Open(filename)
  if err != nil {
    return nil, fmt.Errorf("Unable to read %s: %v", filename, err)
  }

  scanner := bufio.NewScanner(f)
  scanner.Split(scanStmts)

  for scanner.Scan() {
    stmts = append(stmts, scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    return nil, fmt.Errorf("Unable to parse input from %s: %v", filename, err)
  }

  return stmts, nil
}

func scanStmts(data []byte, atEOF bool) (advance int, token []byte, err error) {
  // Skip leading spaces.
  start := 0
  for width := 0; start < len(data); start += width {
    var r rune
    r, width = utf8.DecodeRune(data[start:])
    if !unicode.IsSpace(r) {
      break
    }
  }
  if atEOF && len(data) == 0 {
    return 0, nil, nil
  }

  end := start
  // Scan until semicolon, marking end of statement.
  for width, i := 0, start; i < len(data); i += width {
    var r rune
    r, width = utf8.DecodeRune(data[i:])
    if r == ';' {
      return i + width, data[start:i], nil
    } else if !unicode.IsSpace(r) {
      end = i + 1
    }
  }
  // If we're at EOF, we have a final, non-empty, non-terminated statement. Return it.
  if atEOF && len(data) > start {
    return len(data), data[start:end], nil
  }
  // Request more data.
  return 0, nil, nil
}
