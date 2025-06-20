myapp/
├── cmd/
│   └── myapp/            # entrypoint chính
│       └── main.go       # khởi tạo server, router, và engine template
├── internal/             # logic/handler chỉ dùng nội bộ
│   ├── handler/
│   └── service/
├── pkg/                  # thư viện có thể public (nếu cần)
├── views/                # chứa file .jet
│   ├── layouts/
│   │    └── main.jet     # layout chung
│   ├── partials/
│   │    ├── header.jet
│   │    └── footer.jet
│   └── pages/
│        ├── index.jet
│        └── about.jet
├── web/ (tuỳ chọn)       # assets static: css/js/img
├── go.mod
└── go.sum


| Hành động        | Endpoint               | Method   | Ghi chú                             |
| ---------------- | ---------------------- | -------- | ----------------------------------- |
| Tạo task         | `/tasks`               | `POST`   | Gọi `AddTask()`                     |
| Vô hiệu hoá task | `/tasks/:name/disable` | `PATCH`  | Gọi `DisableTaskByName()`           |
| Bật lại task     | `/tasks/:name/enable`  | `PATCH`  | Gọi `EnableTaskByName()` *(nên có)* |
| Xoá task         | `/tasks/:name`         | `DELETE` | Gọi `DeleteTaskByName()`            |

Dưới đây là danh sách **các RESTful API endpoint cần thiết đầy đủ cho một thực thể (resource)** bất kỳ (ví dụ: `User`, `Product`, `Task`,...) trong kiến trúc **Feature-Based Folder Structure**, bao gồm đầy đủ các query cần thiết và có thể mở rộng.

---

### ✅ **Danh sách RESTful Endpoint Cần Thiết (CRUD + Mở rộng)**

| Method | Endpoint       | Mô tả                                                    |
| ------ | -------------- | -------------------------------------------------------- |
| GET    | `/items`       | Lấy danh sách toàn bộ item (có phân trang, lọc, sắp xếp) |
| GET    | `/items/:id`   | Lấy chi tiết một item                                    |
| POST   | `/items`       | Tạo mới một item                                         |
| PUT    | `/items/:id`   | Cập nhật toàn bộ thông tin item                          |
| PATCH  | `/items/:id`   | Cập nhật một phần thông tin của item                     |
| DELETE | `/items/:id`   | Xóa một item                                             |
| DELETE | `/items`       | Xóa nhiều item (dùng query hoặc body để chỉ định)        |
| GET    | `/items/count` | Đếm số lượng item theo điều kiện                         |
| GET    | `/items/stats` | Thống kê (tùy theo logic)                                |
| POST   | `/items/bulk`  | Tạo nhiều item cùng lúc                                  |
| PUT    | `/items/bulk`  | Cập nhật nhiều item                                      |
| PATCH  | `/items/bulk`  | Cập nhật một phần nhiều item                             |

---

### 📌 **Query Parameters Thường Dùng Cho Endpoint GET `/items`**

```http
GET /items?page=1&limit=20&sort=-createdAt&status=active&q=searchTerm
```

| Tham số     | Mô tả                                                    |
| ----------- | -------------------------------------------------------- |
| `page`      | Trang hiện tại (pagination)                              |
| `limit`     | Số lượng item mỗi trang                                  |
| `sort`      | Sắp xếp theo field (ví dụ: `-createdAt`, `name`)         |
| `q`         | Tìm kiếm toàn cục (full-text search)                     |
| `status`    | Lọc theo trạng thái (`active`, `inactive`,...)           |
| `from`/`to` | Lọc theo khoảng thời gian (`createdAt`, `updatedAt`,...) |
| ...         | Tùy ý mở rộng các trường filter khác theo business       |

---

### 🧩 **Cấu Trúc Thư Mục Feature-Based Cho `items`**

```plaintext
src/
├── items/
│   ├── controller/
│   │   └── items.controller.ts
│   ├── service/
│   │   └── items.service.ts
│   ├── dto/
│   │   ├── create-item.dto.ts
│   │   ├── update-item.dto.ts
│   │   └── query-item.dto.ts
│   ├── entity/
│   │   └── item.entity.ts
│   ├── items.module.ts
│   └── items.route.ts (nếu dùng express/router-style)
```

---

### ✨ **Gợi ý Mở Rộng**

* `GET /items/export`: Export dữ liệu (CSV, Excel)
* `POST /items/import`: Import dữ liệu
* `POST /items/:id/clone`: Clone bản ghi
* `GET /items/:id/history`: Lịch sử chỉnh sửa
* `PATCH /items/:id/restore`: Khôi phục bản ghi đã xoá mềm (nếu dùng soft delete)
* `POST /items/:id/activate` / `deactivate`: Toggle trạng thái

---

### 🧠 Kết luận

Nếu bạn xây hệ thống **API chuẩn, mở rộng và dễ maintain**, thì **Feature-Based Folder Structure + đầy đủ endpoint + query logic theo chuẩn** là bắt buộc. Không cần “RESTful thần thánh hóa”, nhưng đừng làm nửa vời.

Nếu bạn cần ví dụ cụ thể cho 1 resource (VD: `Task`, `User`, `Post`,...), cứ nêu rõ – tôi sẽ viết full luôn.

Bạn muốn tôi tạo sẵn mẫu code (NestJS, Express, Golang...) cho 1 module cụ thể nào không?
