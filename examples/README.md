# Examples

## 📗 Configuration (JSON)

- **Laps**        - Amount of laps for main distance
- **LapLen**      - Length of each main lap
- **PenaltyLen**  - Length of each penalty lap
- **FiringLines** - Number of firing lines per lap
- **Start**       - Planned start time for the first competitor
- **StartDelta**  - Planned interval between starts

## 🏅 Events

All events are characterized by time and event identifier. Outgoing events are events created during program operation. Events related to the "incoming" category cannot be generated and are output in the same form as they were submitted in the input file.

- All events occur sequentially in time: (**Time of event N+1**) $\ge$ (**Time of event N**).
- Time format **HH:MM:SS.sss**.

### 📝 Events format

[***time***] **eventID** **competitorID** extraParams

```ignorelang
Incoming events
EventID | extraParams | Comments
1       |             | The competitor registered
2       | startTime   | The start time was set by a draw
3       |             | The competitor is on the start line
4       |             | The competitor has started
5       | firingRange | The competitor is on the firing range
6       | target      | The target has been hit
7       |             | The competitor left the firing range
8       |             | The competitor entered the penalty laps
9       |             | The competitor left the penalty laps
10      |             | The competitor ended the main lap
11      | comment     | The competitor can`t continue
```

An competitor is disqualified if he/she does not start during his/her start interval. This marked as **NotStarted** in final report.

If the competitor can`t continue it should be marked in final report as **NotFinished**

```ignorelang
Outgoing events
EventID | extraParams | Comments
32      |             | The competitor is disqualified
33      |             | The competitor has finished
```

## 🗒️ Final report

The final report should contain the list of all registered competitors
sorted by ascending time.

- Total time includes the difference between scheduled and actual start time or **NotStarted**/**NotFinished** marks
- Time taken to complete each lap
- Average speed for each lap [m/s]
- Time taken to complete penalty laps
- Average speed over penalty laps [m/s]
- Number of hits/number of shots

## 🔵 Examples

### Single competitor

Run with:

```bash
CONFIG_PATH="examples/single/config.json" go run . < examples/single/events
```

[single/config.json](/examples/single/config.json)

```json
{
    "laps" : 2,
    "lapLen": 3651,
    "penaltyLen": 50,
    "firingLines": 1,
    "start": "09:30:00",
    "startDelta": "00:00:30"
}
```

[examples/single/events](/examples/single/events)

```ignorelang
[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1
[09:49:31.659] 5 1 1
[09:49:33.123] 6 1 1
[09:49:34.650] 6 1 2
[09:49:35.937] 6 1 4
[09:49:37.364] 6 1 5
[09:49:38.339] 7 1
[09:49:55.915] 8 1
[09:51:48.391] 9 1
[09:59:03.872] 10 1
[09:59:03.872] 11 1 Lost in the forest

```

`Output log`:

```ignorelang
[09:05:59.867] The competitor(1) registered
[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000
[09:29:45.734] The competitor(1) is on the start line
[09:30:01.005] The competitor(1) has started
[09:49:31.659] The competitor(1) is on the firing range(1)
[09:49:33.123] The target(1) has been hit by competitor(1)
[09:49:34.650] The target(2) has been hit by competitor(1)
[09:49:35.937] The target(4) has been hit by competitor(1)
[09:49:37.364] The target(5) has been hit by competitor(1)
[09:49:38.339] The competitor(1) left the firing range
[09:49:55.915] The competitor(1) entered the penalty laps
[09:51:48.391] The competitor(1) left the penalty laps
[09:59:03.872] The competitor(1) ended the main lap
[09:59:05.321] The competitor(1) can`t continue: Lost in the forest
[NotFinished] 1 [{00:29:03.872, 2.094}, {,}] {00:01:52.476, 0.445} 4/5
```

### Multiple competitors

Run with:

```bash
CONFIG_PATH="examples/multiple/config.json" go run . < examples/multiple/events
```

[mulitple/config.json](/examples/multiple/config.json)

```json
{
    "laps": 2,
    "lapLen": 3500,
    "penaltyLen": 150,
    "firingLines": 2,
    "start": "10:00:00.000",
    "startDelta": "00:01:30"
}
```

[examples/multiple/events](/examples/multiple/events)

```ignorelang
[09:31:49.285] 1 3
[09:32:17.531] 1 2
[09:37:47.892] 1 5
[09:38:28.673] 1 1
[09:39:25.079] 1 4
[09:55:00.000] 2 1 10:00:00.000
[09:56:30.000] 2 2 10:01:30.000
[09:58:00.000] 2 3 10:03:00.000
[09:59:30.000] 2 4 10:04:30.000
[09:59:45.000] 3 1
[10:00:01.744] 4 1
[10:01:00.000] 2 5 10:06:00.000
[10:01:09.000] 3 2
[10:01:31.503] 4 2
[10:02:36.000] 3 3
[10:03:00.887] 4 3
[10:04:08.000] 3 4
[10:04:31.278] 4 4
[10:05:42.000] 3 5
[10:06:00.331] 4 5
[10:08:49.289] 5 1 1
[10:08:50.884] 6 1 1
[10:08:51.400] 6 1 2
[10:08:52.797] 6 1 5
[10:08:55.658] 7 1
[10:09:03.232] 8 1
[10:10:22.273] 5 2 1
[10:10:23.804] 6 2 1
[10:10:25.036] 6 2 3
[10:10:25.449] 6 2 4
[10:10:26.002] 6 2 5
[10:10:29.125] 7 2
[10:10:38.142] 8 2
[10:10:43.232] 9 1
[10:11:28.142] 9 2
[10:11:54.557] 5 3 1
[10:11:56.076] 6 3 1
[10:11:56.760] 6 3 2
[10:11:57.217] 6 3 3
[10:11:57.659] 6 3 4
[10:11:58.179] 6 3 5
[10:12:01.341] 7 3
[10:12:35.380] 10 1
[10:13:27.246] 5 4 1
[10:13:29.773] 6 4 3
[10:13:30.443] 6 4 4
[10:13:30.836] 6 4 5
[10:13:33.970] 7 4
[10:13:43.912] 8 4
[10:14:09.746] 10 2
[10:15:20.988] 5 5 1
[10:15:22.758] 6 5 1
[10:15:23.083] 6 5 2
[10:15:23.682] 6 5 3
[10:15:23.912] 9 4
[10:15:27.197] 7 5
[10:15:31.757] 8 5
[10:15:43.273] 10 3
[10:17:11.757] 9 5
[10:17:16.947] 10 4
[10:19:21.270] 10 5
[10:21:34.847] 5 1 2
[10:21:36.495] 6 1 1
[10:21:36.920] 6 1 2
[10:21:37.626] 6 1 3
[10:21:38.628] 6 1 5
[10:21:41.449] 7 1
[10:21:50.476] 8 1
[10:22:40.476] 9 1
[10:23:00.773] 5 2 2
[10:23:02.498] 6 2 1
[10:23:02.841] 6 2 2
[10:23:03.453] 6 2 3
[10:23:04.051] 6 2 4
[10:23:07.554] 7 2
[10:23:10.987] 8 2
[10:24:00.987] 9 2
[10:24:43.323] 5 3 2
[10:24:44.954] 6 3 1
[10:24:45.508] 6 3 2
[10:24:45.923] 6 3 3
[10:24:46.559] 6 3 4
[10:24:46.958] 6 3 5
[10:24:49.905] 7 3
[10:25:26.047] 10 1
[10:26:36.573] 5 4 2
[10:26:38.368] 6 4 1
[10:26:38.786] 6 4 2
[10:26:39.113] 6 4 3
[10:26:39.629] 6 4 4
[10:26:40.238] 6 4 5
[10:26:43.208] 7 4
[10:26:48.356] 10 2
[10:28:28.112] 5 5 2
[10:28:29.629] 6 5 1
[10:28:30.408] 6 5 2
[10:28:30.769] 6 5 3
[10:28:31.882] 6 5 5
[10:28:34.274] 7 5
[10:28:34.773] 10 3
[10:28:38.151] 8 5
[10:29:28.151] 9 5
[10:30:36.413] 10 4
[10:32:22.472] 10 5
```

`Output log`:

```ignorelang
[09:31:49.285] The competitor(3) registered
[09:32:17.531] The competitor(2) registered
[09:37:47.892] The competitor(5) registered
[09:38:28.673] The competitor(1) registered
[09:39:25.079] The competitor(4) registered
[09:55:00.000] The start time for the competitor(1) was set by a draw to 10:00:00.000
[09:56:30.000] The start time for the competitor(2) was set by a draw to 10:01:30.000
[09:58:00.000] The start time for the competitor(3) was set by a draw to 10:03:00.000
[09:59:30.000] The start time for the competitor(4) was set by a draw to 10:04:30.000
[09:59:45.000] The competitor(1) is on the start line
[10:00:01.744] The competitor(1) has started
[10:01:00.000] The start time for the competitor(5) was set by a draw to 10:06:00.000
[10:01:09.000] The competitor(2) is on the start line
[10:01:31.503] The competitor(2) has started
[10:02:36.000] The competitor(3) is on the start line
[10:03:00.887] The competitor(3) has started
[10:04:08.000] The competitor(4) is on the start line
[10:04:31.278] The competitor(4) has started
[10:05:42.000] The competitor(5) is on the start line
[10:06:00.331] The competitor(5) has started
[10:08:49.289] The competitor(1) is on the firing range(1)
[10:08:50.884] The target(1) has been hit by competitor(1)
[10:08:51.400] The target(2) has been hit by competitor(1)
[10:08:52.797] The target(5) has been hit by competitor(1)
[10:08:55.658] The competitor(1) left the firing range
[10:09:03.232] The competitor(1) entered the penalty laps
[10:10:22.273] The competitor(2) is on the firing range(1)
[10:10:23.804] The target(1) has been hit by competitor(2)
[10:10:25.036] The target(3) has been hit by competitor(2)
[10:10:25.449] The target(4) has been hit by competitor(2)
[10:10:26.002] The target(5) has been hit by competitor(2)
[10:10:29.125] The competitor(2) left the firing range
[10:10:38.142] The competitor(2) entered the penalty laps
[10:10:43.232] The competitor(1) left the penalty laps
[10:11:28.142] The competitor(2) left the penalty laps
[10:11:54.557] The competitor(3) is on the firing range(1)
[10:11:56.076] The target(1) has been hit by competitor(3)
[10:11:56.760] The target(2) has been hit by competitor(3)
[10:11:57.217] The target(3) has been hit by competitor(3)
[10:11:57.659] The target(4) has been hit by competitor(3)
[10:11:58.179] The target(5) has been hit by competitor(3)
[10:12:01.341] The competitor(3) left the firing range
[10:12:35.380] The competitor(1) ended the main lap
[10:13:27.246] The competitor(4) is on the firing range(1)
[10:13:29.773] The target(3) has been hit by competitor(4)
[10:13:30.443] The target(4) has been hit by competitor(4)
[10:13:30.836] The target(5) has been hit by competitor(4)
[10:13:33.970] The competitor(4) left the firing range
[10:13:43.912] The competitor(4) entered the penalty laps
[10:14:09.746] The competitor(2) ended the main lap
[10:15:20.988] The competitor(5) is on the firing range(1)
[10:15:22.758] The target(1) has been hit by competitor(5)
[10:15:23.083] The target(2) has been hit by competitor(5)
[10:15:23.682] The target(3) has been hit by competitor(5)
[10:15:23.912] The competitor(4) left the penalty laps
[10:15:27.197] The competitor(5) left the firing range
[10:15:31.757] The competitor(5) entered the penalty laps
[10:15:43.273] The competitor(3) ended the main lap
[10:17:11.757] The competitor(5) left the penalty laps
[10:17:16.947] The competitor(4) ended the main lap
[10:19:21.270] The competitor(5) ended the main lap
[10:21:34.847] The competitor(1) is on the firing range(2)
[10:21:36.495] The target(1) has been hit by competitor(1)
[10:21:36.920] The target(2) has been hit by competitor(1)
[10:21:37.626] The target(3) has been hit by competitor(1)
[10:21:38.628] The target(5) has been hit by competitor(1)
[10:21:41.449] The competitor(1) left the firing range
[10:21:50.476] The competitor(1) entered the penalty laps
[10:22:40.476] The competitor(1) left the penalty laps
[10:23:00.773] The competitor(2) is on the firing range(2)
[10:23:02.498] The target(1) has been hit by competitor(2)
[10:23:02.841] The target(2) has been hit by competitor(2)
[10:23:03.453] The target(3) has been hit by competitor(2)
[10:23:04.051] The target(4) has been hit by competitor(2)
[10:23:07.554] The competitor(2) left the firing range
[10:23:10.987] The competitor(2) entered the penalty laps
[10:24:00.987] The competitor(2) left the penalty laps
[10:24:43.323] The competitor(3) is on the firing range(2)
[10:24:44.954] The target(1) has been hit by competitor(3)
[10:24:45.508] The target(2) has been hit by competitor(3)
[10:24:45.923] The target(3) has been hit by competitor(3)
[10:24:46.559] The target(4) has been hit by competitor(3)
[10:24:46.958] The target(5) has been hit by competitor(3)
[10:24:49.905] The competitor(3) left the firing range
[10:25:26.047] The competitor(1) ended the main lap
[10:25:26.047] The competitor(1) has finished
[10:26:36.573] The competitor(4) is on the firing range(2)
[10:26:38.368] The target(1) has been hit by competitor(4)
[10:26:38.786] The target(2) has been hit by competitor(4)
[10:26:39.113] The target(3) has been hit by competitor(4)
[10:26:39.629] The target(4) has been hit by competitor(4)
[10:26:40.238] The target(5) has been hit by competitor(4)
[10:26:43.208] The competitor(4) left the firing range
[10:26:48.356] The competitor(2) ended the main lap
[10:26:48.356] The competitor(2) has finished
[10:28:28.112] The competitor(5) is on the firing range(2)
[10:28:29.629] The target(1) has been hit by competitor(5)
[10:28:30.408] The target(2) has been hit by competitor(5)
[10:28:30.769] The target(3) has been hit by competitor(5)
[10:28:31.882] The target(5) has been hit by competitor(5)
[10:28:34.274] The competitor(5) left the firing range
[10:28:34.773] The competitor(3) ended the main lap
[10:28:34.773] The competitor(3) has finished
[10:28:38.151] The competitor(5) entered the penalty laps
[10:29:28.151] The competitor(5) left the penalty laps
[10:30:36.413] The competitor(4) ended the main lap
[10:30:36.413] The competitor(4) has finished
[10:32:22.472] The competitor(5) ended the main lap
[10:32:22.472] The competitor(5) has finished
[00:25:18.356] 2 [{00:12:39.746, 4.607}, {00:12:38.610, 4.614}] {00:01:40.000, 3.000} 8/10
[00:25:26.047] 1 [{00:12:35.380, 4.633}, {00:12:50.667, 4.542}] {00:02:30.000, 3.000} 7/10
[00:25:34.773] 3 [{00:12:43.273, 4.586}, {00:12:51.500, 4.537}] {,} 10/10
[00:26:06.413] 4 [{00:12:46.947, 4.564}, {00:13:19.466, 4.378}] {00:01:40.000, 3.000} 8/10
[00:26:22.472] 5 [{00:13:21.270, 4.368}, {00:13:01.202, 4.480}] {00:02:30.000, 3.000} 7/10
```
