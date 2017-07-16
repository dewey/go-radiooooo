# Radiooooo

This is a archiving scraper for http://radiooooo.com/. It'll download everything it can find
and stores it into a simple directory structure like that:

```
.
├── ARG
│   └── 1990
│       ├── 11502f64-c30a-430e-b2d9-5543e229a234
│       │   ├── 11502f64-c30a-430e-b2d9-5543e229a234.json
│       │   └── 11502f64-c30a-430e-b2d9-5543e229a234.mp3
│       ├── 8c72651b-9bde-4f57-bdc5-00a4ceb6482f
│       │   └── 8c72651b-9bde-4f57-bdc5-00a4ceb6482f.mp3
│       ├── 98348989-79a2-4a1f-ae77-63e1d27593d3
│       │   └── 98348989-79a2-4a1f-ae77-63e1d27593d3.mp3
│       └── 9fc71867-e091-4899-8e73-f32baba4e647
│           ├── 9fc71867-e091-4899-8e73-f32baba4e647.json
│           └── 9fc71867-e091-4899-8e73-f32baba4e647.mp3
├── ARM
└── BGR
```

It currently downloads the mp3 file only and also archives the raw JSON response for further
processing later on.