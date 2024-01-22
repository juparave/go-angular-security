export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string;

  // transient members
  token: string;
  refreshToken: string;
}
