// ref: https://bobbyhadz.com/blog/javascript-initialize-date-with-timezone
export function changeTimeZone(date: Date | string, timeZone: string) {
  if (typeof date === 'string') {
    return new Date(
      new Date(date).toLocaleString('es-MX', {
        timeZone,
      }),
    );
  }

  return new Date(
    date.toLocaleString('es-MX', {
      timeZone,
    }),
  );
}