import { Pipe, PipeTransform } from '@angular/core';

/**
 * The `role` pipe checks if a given role name is contained in an array of roles.
 *
 * @example
 * // In your component template:
 * // Assuming user?.roles is an array of roles
 * <span *ngIf="user?.roles | role:'admin'">Administrator</span>
 *
 * @param {string[]} roles - The array of roles to check against.
 * @param {string} roleName - The role name to check for in the array.
 * @returns {boolean} - Returns `true` if the role is present, otherwise `false`.
 */
@Pipe({
  name: 'role',
  standalone: true
})
export class RolePipe implements PipeTransform {

  transform(roles: string | string[] | undefined, roleName: string): boolean {
    if (!roles) {
      return false;
    }
    if (typeof roles === 'string') {
      roles = roles.split(',');
    }
    return roles.includes(roleName);
  }

}
