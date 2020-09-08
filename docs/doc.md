## 说明

由于某站电费格式错乱，解析复杂，所以对目前的学生宿舍楼栋进行了一个筛选和整理：

+ 西区1-8栋
+ 东区1-16，东13栋东/西，东15东/西、东附一
+ 元宝山1-5栋
+ 南湖1-13东
+ 国交（重点）：4/8/9 栋

国交目前只选择国4、国8和国9三所楼栋，确定这三栋是目标学生群体。如之后另有情况再进行调整。

## API
### 查询电费

| Method | URL         | Header  |
| ------ | ----------- | ------- |
| GET    | /api/ele/v2 | null    |

**Request Params:**
```
    building: string // 宿舍楼栋简写，如：东16、东15-东。
    room: string // 宿舍号，如：101、101A
```

**RESPONSE Data:**
```json
{
    "code": 0,
    "message": "OK",
    "data": {
        "building": "东7",
        "room": "101A",
        // 是否有照明电费的信息
        // 有些宿舍目前没有（或是暂时找不到、理不清），所以没有的就显示“暂无照明信息即可”
        "has_light": true,
        "light": {
            "kind": "light",
            "remain_power": "228.89",
            "read_time": "2020-09-08 02:04:45",
            "consumption": {
                "usage": "2.56",
                "charge": "1.52"
            }
        },
        "has_air": true,
        "air": {
            "kind": "air",
            "remain_power": "37.82",
            "read_time": "2020-09-08 02:11:07",
            "consumption": {
                "usage": "5.96",
                "charge": "3.55"
            }
        },
        // 除了照明和空调以外，可能还有其它的电费需要负责。
        // 比如东7宿舍还有客厅这一个功能区，因而也要负责客厅的空调这一部分。
        // （目前就一个东7会有这部分的数据）
        "has_more": true,
        "more_data": [
            {
                "kind": "客厅-空调",
                "remain_power": "14.33",
                "read_time": "2020-09-08 02:06:34",
                "consumption": {
                    "usage": "4.15",
                    "charge": "2.47"
                }
            }
        ]
    }
}
```

### 查询宿舍号

| Method | URL               | Header  |
| ------ | ----------------- | ------- |
| GET    | /api/ele/v2/dorms | null    |

**Request Params:**
```
    building: string // 楼栋名，如：南湖1
```

**RESPONSE Data:**
```json
{
    "code": 0,
    "message": "OK",
    "data": {
        "count": 4,
        "list": ["101","102","103","104"]
    }
}
```

### 查询楼栋

| Method | URL                   | Header  |
| ------ | --------------------- | ------- |
| GET    | /api/ele/v2/buildings | null    |

**Request Params:**
```
    area: string // 区域，西区/东区/元宝山/南湖/国交
```

**RESPONSE Data:**
```json
{
    "code": 0,
    "message": "OK",
    "data": [
        {
            "name": "元宝山1栋",
            "alias": "元1"
        },
        {
            "name": "元宝山2栋",
            "alias": "元2"
        },
        {
            "name": "元宝山3栋",
            "alias": "元3"
        },
        {
            "name": "元宝山4栋",
            "alias": "元4"
        },
        {
            "name": "元宝山5栋",
            "alias": "元5"
        }
    ]
}
```

PS：`name` 用于用户展示，`alias`用于请求宿舍号和电费
