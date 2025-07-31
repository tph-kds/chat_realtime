

// import type { FormatMessageTimeProps } from "../configs/types";


function formatMessageTime(date: string | number | Date): string {
  return new Date(date).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
}

export default formatMessageTime;