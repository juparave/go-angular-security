# Angular Frontend for Secure Web Application

This project is the Angular frontend for a secure web application. It interacts with a Go backend to provide features like user authentication, authorization, and management of resources. It was generated with [Angular CLI](https://github.com/angular/angular-cli).

## Project Overview

The frontend is designed to demonstrate best practices in building secure and modular Angular applications. It includes:
- User registration, login, and profile management.
- Role-based access control to different parts of the application.
- Secure handling of JWTs for session management.
- Interaction with a RESTful API provided by the Go backend.

## Development Server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if you change any of the source files.

The backend API is expected to be running on `http://localhost:3000`. A proxy configuration is in place (`src/proxy.config.json`) to forward API requests from `http://localhost:4200/api` to `http://localhost:3000/api`.

## Code Scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use `ng build --prod` for a production build.

## Running Unit Tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running End-to-End Tests

End-to-end tests are not currently set up in this project. If you wish to add them, you would typically use a library like Cypress or Protractor. After setting up, you would run `ng e2e` to execute the end-to-end tests. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Key Modules

The application is structured into several key Angular modules:

- **`AppModule` (`src/app/app.module.ts`):** The root module that bootstraps the application.
- **`AppRoutingModule` (`src/app/app-routing.module.ts`):** Defines the main application routes.
- **`AuthModule` (`src/app/modules/auth/auth.module.ts`):** Handles authentication-related components and services (login, registration).
- **`MainModule` (`src/app/modules/main/main.module.ts`):** Contains the core features of the application accessible after login (e.g., dashboard, user management).
- **`SharedModule` (`src/app/shared/shared.module.ts`):** Contains reusable components, directives, and pipes used across different modules.

## Key Components

Some of the important components include:

- **`AppComponent` (`src/app/app.component.ts`):** The main application shell.
- **`LoginComponent` (`src/app/modules/auth/components/login/login.component.ts`):** Handles user login.
- **`RegisterComponent` (`src/app/modules/auth/components/register/register.component.ts`):** Handles user registration.
- **`DashboardComponent` (`src/app/modules/main/components/dashboard/dashboard.component.ts`):** The main view after a user logs in.
- **`UsersComponent` (`src/app/modules/main/components/users/users.component.ts`):** Component for listing and managing users (example, might require admin privileges).
- **`NavComponent` (`src/app/shared/components/nav/nav.component.ts`):** The main navigation bar.

## Key Services

Core services facilitating application logic:

- **`AuthService` (`src/app/modules/auth/services/auth.service.ts`):** Manages authentication logic, user sessions, and API calls to auth endpoints.
- **`UserService` (`src/app/services/user.service.ts`):** (Example) Manages user-related data and API calls.
- **`RoleService` (`src/app/services/role.service.ts`):** (Example) Manages role-related data.
- **`AuthGuard` (`src/app/modules/auth/guards/auth.guard.ts`):** A route guard that prevents access to certain routes if the user is not authenticated.

## Styling with Tailwind CSS

This project uses [Tailwind CSS](https://tailwindcss.com/) for styling.
To set it up (if not already done):
1. Install Tailwind CSS:
   ```bash
   npm install -D tailwindcss postcss autoprefixer
   npx tailwindcss init
   ```
2. Configure `tailwind.config.js` to include paths to your HTML and TypeScript files:
   ```js
   module.exports = {
     content: [
       "./src/**/*.{html,ts}",
     ],
     theme: {
       extend: {},
     },
     plugins: [],
   }
   ```
3. Include Tailwind directives in `src/styles.scss`:
   ```scss
   @tailwind base;
   @tailwind components;
   @tailwind utilities;
   ```

## Further Help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
