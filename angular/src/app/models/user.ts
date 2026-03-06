import { Subscription } from "./subscription";

export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string;
  subscription: Subscription

  // transient members
  accessToken: string;
  refreshToken: string;
}
