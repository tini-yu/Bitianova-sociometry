import "./matrix.css";
import {
    GetResults,
    LoadResult,
    SaveResult,
    CheckMatrices
} from "../wailsjs/go/main/App";

const COLUMN_LABELS = [
    "Положительные выборы",
    "Отрицательные выборы",
    "Статус",
    "Кол-во взаимных социометрических выборов (положительных)",
    "Кол-во взаимных социометрических выборов (отрицательных)",
    "Кол-во противоречивых выборов",
    "Аутосоциометрия",
    "Кол-во референтных выборов",
    "Целевая группа" 
] as const;

const EDITABLE_COLUMN_INDEX = 8;

const EDIT_OPTIONS = ["-", "+"] as const;
type EditOption = typeof EDIT_OPTIONS[number];

class Matrix4 {
    private labels: string[] = [];
    private matrix: string[][] = [];

    private currentUUID = "";
    private filename = "results.json";

    private tableHead!: HTMLTableSectionElement;
    private tableBody!: HTMLTableSectionElement;
    private statusEl!: HTMLElement;

    private autoMaps: Record<string, string>[] = [];

    constructor() {
        this.waitForDOM();
    }

    private async waitForDOM() {
        await this.waitForElement("#matrix-table");
        await this.waitForElement("#matrix-status");

        this.tableHead = document.querySelector("#matrix-table thead")!;
        this.tableBody = document.querySelector("#matrix-table tbody")!;
        this.statusEl = document.getElementById("matrix-status")!;

        document.getElementById("save-btn")!.onclick = () => this.save();
        document.getElementById("load-btn")!.onclick = () => this.load();

        await this.load();
    }

    private waitForElement(selector: string): Promise<void> {
        return new Promise(resolve => {
            if (document.querySelector(selector)) return resolve();
            const obs = new MutationObserver(() => {
                if (document.querySelector(selector)) {
                    obs.disconnect();
                    resolve();
                }
            });
            obs.observe(document.body, { childList: true, subtree: true });
        });
    }

    private async load() {
        this.showStatus("Идёт загрузка...", "blue");

        const ok = await CheckMatrices();
        if (!ok) {
            this.tableHead.innerHTML = "";
            this.tableBody.innerHTML = "";
            this.showStatus("Необходимо заполнить матрицы вопросов", "red");
            return;
        }

        const { labels, maps, uuid } = await GetResults();

        this.labels = labels;
        this.currentUUID = uuid;
        this.autoMaps = maps;

        const rows = labels.length;
        const cols = COLUMN_LABELS.length;
        //на всякий случай дефолтная пустая матрица
        this.matrix = Array.from({ length: rows }, () =>
            Array.from({ length: cols }, () => "-")
        );

        // автозаполнение колонок 1-8
        for (let r = 0; r < rows; r++) {
            const label = this.labels[r];
            for (let c = 0; c < EDITABLE_COLUMN_INDEX; c++) {
                const map = this.autoMaps[c];
                this.matrix[r][c] = map[label] ?? "-";
            }
        }

        try {
            const json = await LoadResult(this.filename);
            const saved = JSON.parse(json) as { uuid: string; data: string[][] };

            if (
                saved.uuid === this.currentUUID &&
                saved.data.length === rows &&
                saved.data.every(row => row.length === cols)
            ) {
                // все колонки кроме 9 автозаполняются. 9 берём из сохранений
                for (let r = 0; r < rows; r++) {
                    this.matrix[r][EDITABLE_COLUMN_INDEX] =
                        saved.data[r][EDITABLE_COLUMN_INDEX];
                }
                this.showStatus("Данные загружены", "green");
            }
        } catch {
            this.showStatus("Нет данных для загрузки", "gray");
        }

        await this.render();
    }

    private async render() {

        const ok = await CheckMatrices();
        if (!ok) {
            this.tableHead.innerHTML = "";
            this.tableBody.innerHTML = "";
            this.showStatus("Необходимо заполнить матрицы вопросов", "red");
            return;
        }

        try {
            const { labels, maps, uuid } = await GetResults();
            this.labels = labels;
            this.currentUUID = uuid;
            this.autoMaps = maps;

            for (let r = 0; r < labels.length; r++) {
                const label = labels[r];
                for (let c = 0; c < EDITABLE_COLUMN_INDEX; c++) {
                    this.matrix[r][c] = this.autoMaps[c][label] ?? "-";
                }
            }
        } catch {
            this.showStatus("Ошибка рендеринга!", "red");
        }
        const table = document.getElementById("matrix-table") as HTMLTableElement;
        table.innerHTML = "";

        const headerRow = table.insertRow();
        headerRow.classList.add("sticky-header-row");

        const cornerTh = document.createElement("th");
        cornerTh.className = "corner-header";
        headerRow.appendChild(cornerTh);

        COLUMN_LABELS.forEach(label => {
            const th = document.createElement("th");
            th.className = "col-header";
            th.textContent = label;
            headerRow.appendChild(th);
        });

        this.labels.forEach((rowLabel, r) => {
            const tr = table.insertRow();

            const rowTh = document.createElement("th");
            rowTh.className = "row-header";
            rowTh.textContent = rowLabel;
            tr.appendChild(rowTh);

            COLUMN_LABELS.forEach((_, c) => {
                const td = tr.insertCell();

                if (c === EDITABLE_COLUMN_INDEX) {
                    const select = document.createElement("select");
                    EDIT_OPTIONS.forEach(v => {
                        const option = document.createElement("option");
                        option.value = v;
                        option.textContent = v;
                        select.appendChild(option);
                    });
                    select.value = this.matrix[r][c];
                    select.onchange = e => {
                        this.matrix[r][c] =
                            (e.target as HTMLSelectElement).value as EditOption;
                    };
                    td.appendChild(select);
                } else {
                    // auto-filled cells
                    td.textContent = this.matrix[r][c];
                    td.classList.add("readonly-cell");
                }
            });
        });
    }

    private async save() {
        const rows = this.labels.length;
        const data = Array.from({ length: rows }, (_, r) => [...this.matrix[r]]);

        try {
            await SaveResult(this.filename, data, this.currentUUID);
            this.showStatus("Сохранено!", "green");
        } catch (err: any) {
            this.showStatus(err, "red");
        }
    }

    private showStatus(msg: string, color: string) {
        this.statusEl.textContent = msg;
        this.statusEl.style.color = color;
        this.statusEl.style.fontSize = "200%";
        setTimeout(() => {
            if (this.statusEl.textContent === msg) this.statusEl.textContent = "";
        }, 3000);
    }
}

export function init() {
    new Matrix4(); 
}