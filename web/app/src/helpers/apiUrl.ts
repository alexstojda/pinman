export default function apiUrl(path: string) {
  if (
    process.env.REACT_APP_API_HOST !== undefined &&
    process.env.REACT_APP_API_HOST !== ""
  ) {
    return process.env.REACT_APP_API_HOST + path;
  } else {
    return path;
  }
}
