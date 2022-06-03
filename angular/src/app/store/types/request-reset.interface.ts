export interface RequestResetPassword {
  email: string;
}

export interface ValidateRequestResetPassword {
  email: string;
  code: string;
}

export interface ResetPassword {
  email: string;
  newPassword: string;
  confirmPassword: string;
}
