export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;

  // transient members
  token: string;
  refreshToken: string;
}
