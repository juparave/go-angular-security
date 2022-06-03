export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;

  // transient members
  accessToken: string;
  refreshToken: string;
}
