export interface LoginRequest {
  user: {
    email: string;
    password: string;
  };
}

export type ReturnUrl = string;
