# Design System

UI patterns and conventions for the go-angular-security Angular frontend.

---

## Stack

- **Tailwind CSS v4** — utility-first CSS via `@import 'tailwindcss'`
- **DaisyUI 5** — component library via `@plugin 'daisyui'`
- **Angular Material** — used for form fields and overlays where DaisyUI falls short
- **Inter** — primary font family

---

## Theme Configuration

Themes are declared in `src/tailwind.css` using Tailwind v4 `@theme` directive:

```css
@import 'tailwindcss';
@plugin 'daisyui';

:root {
  --daisyui-themes: cupcake, dark;
}

@theme {
  --default-font-family: 'Inter', ui-sans-serif, system-ui, sans-serif;
  --breakpoint-3xl: 1920px;

  /* Brand color tokens */
  --color-primary:   oklch(0.6 0.2 250);   /* indigo */
  --color-secondary: oklch(0.7 0.15 180);  /* teal */
  --color-accent:    oklch(0.65 0.2 40);   /* amber */
}
```

DaisyUI themes (`cupcake` for light, `dark` for dark mode) are toggled by setting `data-theme` on `<html>`.

---

## Color Tokens

Use DaisyUI semantic tokens in templates. Never hardcode hex or RGB values.

| Token | Usage |
|-------|-------|
| `bg-base-100` | Page background |
| `bg-base-200` | Card / sidebar background |
| `bg-base-300` | Divider / subtle background |
| `text-base-content` | Body text |
| `bg-primary` / `text-primary` | Primary actions |
| `bg-secondary` / `text-secondary` | Secondary actions |
| `bg-accent` / `text-accent` | Highlights |
| `bg-error` / `text-error` | Errors |
| `bg-success` / `text-success` | Success states |
| `bg-warning` / `text-warning` | Warnings |

---

## Component Patterns

### Buttons

```html
<!-- Primary action -->
<button class="btn btn-primary">Save</button>

<!-- Secondary -->
<button class="btn btn-secondary">Cancel</button>

<!-- Destructive -->
<button class="btn btn-error">Delete</button>

<!-- Ghost / icon only -->
<button class="btn btn-ghost btn-sm">
  <span class="material-icons">edit</span>
</button>

<!-- Loading state -->
<button class="btn btn-primary" [disabled]="loading">
  <span *ngIf="loading" class="loading loading-spinner loading-sm"></span>
  {{ loading ? 'Saving…' : 'Save' }}
</button>
```

### Cards

```html
<div class="card bg-base-100 shadow-sm border border-base-200">
  <div class="card-body">
    <h2 class="card-title text-lg">Card Title</h2>
    <p class="text-base-content/70">Supporting text.</p>
    <div class="card-actions justify-end mt-4">
      <button class="btn btn-primary btn-sm">Action</button>
    </div>
  </div>
</div>
```

### Forms

Use DaisyUI form classes; Angular Material `mat-form-field` is acceptable for complex inputs.

```html
<label class="form-control w-full">
  <div class="label">
    <span class="label-text">Email</span>
  </div>
  <input
    type="email"
    class="input input-bordered w-full"
    [class.input-error]="form.get('email')?.invalid && form.get('email')?.touched"
    formControlName="email"
    placeholder="you@example.com"
  />
  <div class="label" *ngIf="form.get('email')?.invalid && form.get('email')?.touched">
    <span class="label-text-alt text-error">Valid email required</span>
  </div>
</label>
```

### Tables

```html
<div class="overflow-x-auto rounded-lg border border-base-200">
  <table class="table table-zebra">
    <thead>
      <tr class="bg-base-200">
        <th>Name</th>
        <th>Role</th>
        <th>Status</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr *ngFor="let member of members">
        <td>{{ member.firstName }} {{ member.lastName }}</td>
        <td><span class="badge badge-ghost">{{ member.role }}</span></td>
        <td>
          <span class="badge" [class]="member.enabled ? 'badge-success' : 'badge-error'">
            {{ member.enabled ? 'Active' : 'Disabled' }}
          </span>
        </td>
        <td class="text-right">
          <button class="btn btn-ghost btn-xs">Edit</button>
        </td>
      </tr>
    </tbody>
  </table>
</div>
```

### Navigation / Sidebar

```html
<ul class="menu menu-vertical gap-1 p-2">
  <li>
    <a routerLink="/dashboard" routerLinkActive="active">
      <span class="material-icons text-lg">dashboard</span>
      Dashboard
    </a>
  </li>
  <li>
    <a routerLink="/team" routerLinkActive="active">
      <span class="material-icons text-lg">group</span>
      Team
    </a>
  </li>
</ul>
```

### Alerts / Toasts

```html
<!-- Inline alert -->
<div class="alert alert-error" *ngIf="error">
  <span class="material-icons">error</span>
  <span>{{ error }}</span>
</div>

<!-- DaisyUI toast (positioned) -->
<div class="toast toast-end">
  <div class="alert alert-success">
    <span>Saved successfully.</span>
  </div>
</div>
```

### Modal / Dialog

```html
<!-- DaisyUI modal -->
<dialog id="confirm_modal" class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Confirm Delete</h3>
    <p class="py-4">This action cannot be undone.</p>
    <div class="modal-action">
      <button class="btn btn-error" (click)="onConfirm()">Delete</button>
      <button class="btn" onclick="confirm_modal.close()">Cancel</button>
    </div>
  </div>
</dialog>
```

---

## Layout Guidelines

### Page shell

```html
<div class="flex min-h-screen bg-base-200">
  <!-- Sidebar -->
  <aside class="w-64 bg-base-100 border-r border-base-200 flex-shrink-0">
    <!-- nav -->
  </aside>

  <!-- Main content -->
  <main class="flex-1 p-6 overflow-auto">
    <div class="max-w-5xl mx-auto space-y-6">
      <!-- page content -->
    </div>
  </main>
</div>
```

### Section header

```html
<div class="flex items-center justify-between mb-6">
  <div>
    <h1 class="text-2xl font-bold text-base-content">Page Title</h1>
    <p class="text-sm text-base-content/60 mt-1">Subtitle or description</p>
  </div>
  <button class="btn btn-primary btn-sm">Add Item</button>
</div>
```

---

## Responsive Breakpoints

| Name | Min-width |
|------|-----------|
| `sm` | 640px |
| `md` | 768px |
| `lg` | 1024px |
| `xl` | 1280px |
| `2xl` | 1536px |
| `3xl` | 1920px (custom) |

Use mobile-first: base classes apply to mobile, breakpoint prefixes for larger screens.

```html
<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
```

---

## Angular Material Integration

Angular Material is available for components not covered by DaisyUI (e.g. date pickers, autocomplete). When using it:

- Apply `mat-form-field` with `appearance="outline"` to match the DaisyUI aesthetic
- Override Material colors with CSS custom properties to match DaisyUI theme tokens
- Prefer DaisyUI for simple inputs, buttons, badges, and cards

---

## Conventions

- Use `@for` and `@if` (Angular 17+ control flow) instead of `*ngFor` / `*ngIf` in new components
- Prefer standalone components
- Keep template logic minimal — move complex expressions to component class
- Do not hardcode colors; always use DaisyUI semantic tokens or `@theme` variables
