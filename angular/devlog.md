# Init project

### Angular CLI
    npm install -g @angular/cli@17

### Create new project named `angular` without the default `standalone` and `cd` into it
    ng new angular --standalone false
    cd angular/

### eslint
    npx eslint --init

### Angular Material
    ng add @angular/material

### NGRX
    ng add @ngrx/store@latest
    ng add @ngrx/effects@latest
    npm install --save-dev @ngrx/store-devtools


# Adding basic modules and components

    ng generate module modules/app
    ng generate component modules/app/base
    ng generate component modules/app/custom-sidenav
    ng generate component modules/app/dashboard
    ng generate pipe shared/pipes/role --skip-tests true

