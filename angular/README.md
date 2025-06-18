# Angular Frontend

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 17.1.0. It provides the user interface for the application and interacts with the Go backend API to deliver a seamless and secure user experience.

## Project Structure

The `src/app` directory contains the core of the Angular application. Key subdirectories include:

*   `model/`: Defines TypeScript interfaces and classes representing data structures used throughout the application (e.g., `user.ts`).
*   `modules/`: Contains feature modules that encapsulate specific functionalities. This project includes modules for general application features (`app/`) and user authentication (`login/`).
*   `services/`: Houses injectable services responsible for tasks like API communication, authentication (`auth/auth.service.ts`, `auth.guard.ts`, `refresh-token.interceptor.ts`), file uploads, notifications, and managing application loading state.
*   `shared/`: Includes reusable components (e.g., `page-not-found`), shared modules (e.g., for alerts, displaying backend error messages), custom pipes (e.g., `role.pipe.ts`), and utility functions.
*   `store/`: Implements state management, likely using NgRx, with subdirectories for `actions`, `effects`, `reducers`, and `selectors` to manage application state in a predictable way, particularly for authentication.

## Development Server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if you change any of the source files.

### Backend Interaction

During development, the Angular frontend makes API calls to the Go backend. Ensure the Go backend server is running (typically on a port like `8080`, as described in the main `README.md`). API calls are usually configured in the environment files (`src/environments/`) or directly within the services that communicate with the backend.

*   **Note:** This project does not currently use a `proxy.config.json` for API request proxying. API calls are made directly to the backend URL.

## Code Scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use `ng build --prod` for a production build.

## Running Unit Tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running End-to-End Tests

Run `ng e2e` to execute the end-to-end tests via a platform of your choice. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Further Help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
