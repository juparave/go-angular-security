# Generic Pagination Guide (`Paginate[T]`)

This document guides the engineering team on how to implement and consume type-safe pagination, sorting, and filtering in the backend and frontend.

---

## 1. Overview & Rationale

### Why We Replaced the Legacy `Entity` Interface
In older projects, pagination relied on every model implementing an `Entity` interface with `Count(db)` and `Take(db, limit, offset) interface{}`.

**Problems with that approach:**
* ❌ Required copy-pasted `Take` and `Count` boilerplate on every model.
* ❌ Returned untyped `interface{}` / `any`.
* ❌ Risk of GORM query state corruption between `Count` and `Find`.

**The Modern Go 1.24+ Generic Approach (`utils.Paginate[T]`):**
* ✅ **Zero Model Boilerplate:** Your GORM models need no pagination methods.
* ✅ **Type Safety:** Returns `PaginateResponse[T]` with strict compile-time types.
* ✅ **100% Angular Compatible:** Emits the exact same JSON format (`data`, `total`, `page`, `page_size`, `last_page`).
* ✅ **SQL Injection Safe:** `sort_col` is validated against safe identifier patterns.

---

## 2. API Contract & Query Parameters

### Request Query Parameters
The backend helper parses both Angular Material conventions and standard web conventions:

| Parameter | Aliases | Default | Description |
| :--- | :--- | :--- | :--- |
| `page_index` | `page`, `pageNum` | `0` | 0-indexed page number |
| `page_size` | `pageSize`, `limit` | `25` | Number of items per page (max `100`) |
| `sort_col` | `sortCol`, `sort` | `""` | Column name to sort by (e.g. `name`, `created_at`) |
| `sort_ord` | `sortOrd`, `order` | `ASC` | `ASC` or `DESC` |
| `filter` | `q` | `""` | Generic search keyword |

### Response JSON Payload
```json
{
  "data": [
    {
      "id": "2d8f9a0b1c2e",
      "name": "Acme Corp",
      "email": "contact@acme.com"
    }
  ],
  "total": 142,
  "page": 0,
  "page_size": 25,
  "last_page": 6
}
```

---

## 3. Backend Usage (Go + Fiber + GORM)

### Basic Handler Example

```go
package handlers

import (
	"net/http"
	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// GetAccounts handles GET /api/v1/accounts
func GetAccounts(c *fiber.Ctx) error {
	// 1. Automatically extract & sanitize query params (page_index, page_size, sort_col, etc.)
	params := utils.ParsePaginateParams(c, 25) // 25 is default page_size

	// 2. Build your GORM query with business filters
	db := database.Manager.GetMasterDB().Model(&models.Account{})

	if params.Filter != "" {
		searchTerm := "%" + params.Filter + "%"
		db = db.Where("name LIKE ? OR contact_email LIKE ?", searchTerm, searchTerm)
	}

	// 3. Execute type-safe pagination in one line
	response, err := utils.Paginate[models.Account](db, params)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch accounts",
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}
```

### Preloading Relational Data
You can chain GORM preloads before passing `db` to `Paginate[T]`:

```go
db := tenantDB.Model(&models.User{}).Preload("Role")
response, err := utils.Paginate[models.User](db, params)
```

---

## 4. Frontend Integration (Angular 19)

### TypeScript Interface
In your Angular models (`src/app/core/models/paginate.model.ts`):

```typescript
export interface PaginateResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  last_page: number;
}

export interface PaginateParams {
  page_index: number;
  page_size: number;
  sort_col?: string;
  sort_ord?: 'ASC' | 'DESC';
  filter?: string;
}
```

### Angular Service
```typescript
@Injectable({ providedIn: 'root' })
export class AccountService {
  private http = inject(HttpClient);
  private apiUrl = '/api/v1/accounts';

  getAccounts(params: PaginateParams): Observable<PaginateResponse<Account>> {
    let httpParams = new HttpParams()
      .set('page_index', params.page_index.toString())
      .set('page_size', params.page_size.toString());

    if (params.sort_col) httpParams = httpParams.set('sort_col', params.sort_col);
    if (params.sort_ord) httpParams = httpParams.set('sort_ord', params.sort_ord);
    if (params.filter) httpParams = httpParams.set('filter', params.filter);

    return this.http.get<PaginateResponse<Account>>(this.apiUrl, { params: httpParams });
  }
}
```

### Angular Component with MatPaginator & MatSort
```typescript
@Component({ ... })
export class AccountsTableComponent implements AfterViewInit {
  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  dataSource = new MatTableDataSource<Account>([]);
  totalRecords = 0;
  isLoading = false;

  private accountService = inject(AccountService);

  ngAfterViewInit() {
    // Reset to page 0 when sorting changes
    this.sort.sortChange.subscribe(() => (this.paginator.pageIndex = 0));

    // Combine sort and page changes into server request
    merge(this.sort.sortChange, this.paginator.page)
      .pipe(
        startWith({}),
        switchMap(() => {
          this.isLoading = true;
          return this.accountService.getAccounts({
            page_index: this.paginator.pageIndex,
            page_size: this.paginator.pageSize,
            sort_col: this.sort.active,
            sort_ord: this.sort.direction.toUpperCase() as 'ASC' | 'DESC',
          }).pipe(catchError(() => of(null)));
        }),
        map(response => {
          this.isLoading = false;
          if (response === null) return [];
          this.totalRecords = response.total;
          return response.data;
        })
      )
      .subscribe(data => (this.dataSource.data = data));
  }
}
```

---

## 5. Security & Best Practices

1. **SQL Injection Defense:** `sort_col` is validated in `utils.ParsePaginateParams` against `^[a-zA-Z0-9_\.]+$`. Any parameter containing semicolons, comments, or SQL keywords is automatically stripped.
2. **Page Size Caps:** Maximum `page_size` is hard-capped at `100` to prevent memory exhaustion from oversized queries.
3. **Session Isolation:** `Paginate[T]` clones the GORM session for the `Count(&total)` query. This guarantees that `Limit()`, `Offset()`, or `Order()` do not corrupt the count calculation.
