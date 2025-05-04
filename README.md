# goathlon

![LICENSE-shield][license-shield-url] ![WORKFLOW-status][workflow-status-url]

A prototype system for processing and analyzing biathlon competition events.

## 🔧 Installation & Setup

### Prerequisites

- **Go**: v1.20 or newer.

### Usage

You need to complete a few steps in order to run the program:

1. `CONFIG_PATH` environment varaible must be set.
2. Input stream must be provided.

For example:

```bash
CONFIG_PATH="test/testdata/single/config.json" go run . < test/testdata/single/events
```

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

`test/testdata/single/config.json`

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

`test/testdata/single/events`

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

Run with:

```bash
CONFIG_PATH="test/testdata/single/config.json" go run . < test/testdata/single/events
```

Output log:

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

<!---->
[license-shield-url]: https://img.shields.io/github/license/artem-burashnikov/goathlon?style=for-the-badge&color=blue
[workflow-status-url]: https://img.shields.io/github/actions/workflow/status/artem-burashnikov/goathlon/.github%2Fworkflows%2Fci.yaml?style=for-the-badge&color=lightgreen
