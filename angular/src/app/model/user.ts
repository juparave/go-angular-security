export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string;

  // transient members
  accessToken: string;
  refreshToken: string;
}
