import "./matrix.css";
import { ListExport, LoadMatrix, SaveMatrix } from "../wailsjs/go/main/App";

const OPTIONS = ["-", "Хочу видеть"] as const;
type Option = typeof OPTIONS[number];

const OPTIONS_MAP = {
    0: "-",
    1: "Хочу видеть",
} as const;

const REVERSE_MAP = {
    "-": 0,
    "Хочу видеть": 1,
} as const;

class MatrixApp {
    private labels: string[] = [];
    private currentUUID = "";
    private size = 0;
    private matrix: Option[][] = [];
    private statusEl!: HTMLElement;
    private matrixFilename = "matrix3.json";

    constructor() {
        this.waitForDOMAndInit();
    }


    private async waitForDOMAndInit() {

        await this.waitForElement("#matrix-table");
        await this.waitForElement("#matrix-status");
        await this.waitForElement("#save-btn");

        this.statusEl = document.getElementById("matrix-status")!;

        document.getElementById("save-btn")!.onclick = () => this.saveMatrix();
        document.getElementById("load-btn")!.onclick = () => this.refreshTable();

        await this.refreshTable();
    }

    private waitForElement(selector: string, timeout = 5000): Promise<void> {
        return new Promise((resolve, reject) => {
            if (document.querySelector(selector)) {
                return resolve();
            }

            const observer = new MutationObserver(() => {
                if (document.querySelector(selector)) {
                    observer.disconnect();
                    resolve();
                }
            });

            observer.observe(document.body, {
                childList: true,
                subtree: true,
            });

            setTimeout(() => {
                observer.disconnect();
                reject(new Error(`Timeout waiting for ${selector}`));
            }, timeout);
        });
    }

    private async refreshTable() {
        this.showStatus("Идёт загрузка...", "blue");

        try {
            const result = await ListExport(); // json: labels / uuid

            this.labels = result.labels; //имена
            this.currentUUID = result.uuid;

            console.log("UUid", this.currentUUID);
            this.size = this.labels.length;
            this.matrix = Array(this.size)
                .fill(null)
                .map(() => Array(this.size).fill("-" as Option));

            try {
                const savedJson: string = await LoadMatrix(this.matrixFilename);
                const savedData = JSON.parse(savedJson) as { uuid: string; data: number[][] };

                if (savedData.uuid === this.currentUUID &&
                    savedData.data.length === this.size &&
                    savedData.data.every(row => row.length === this.size)) {
                    this.matrix = savedData.data.map(row =>
                        row.map(n => OPTIONS_MAP[n as 0 | 1])
                    );
                    this.showStatus("Данные загружены", "green");
                } else {
                    this.showStatus("Загрузка невозможна - список имён был изменён", "orange");
                }
            } catch (e) {
                this.showStatus("Нет данных для загрузки", "gray");
            }


            this.renderTable();
        } catch (err: any) {
            this.showStatus(`Error: ${err.message || err}`, "red");
            console.error(err);
        }
    }

    private renderTable() {
        const table = document.getElementById("matrix-table") as HTMLTableElement;
        table.innerHTML = "";

        const headerRow = table.insertRow();
        headerRow.classList.add("sticky-header-row");

        const cornerTh = document.createElement("th");
        cornerTh.className = "corner-header";
        headerRow.appendChild(cornerTh);

        this.labels.forEach(label => {
            const th = document.createElement("th");
            th.className = "col-header";
            th.textContent = label;
            headerRow.appendChild(th);
        });

        this.labels.forEach((rowLabel, rowIdx) => {
            const tr = table.insertRow();

            const rowTh = document.createElement("th");
            rowTh.className = "row-header";
            rowTh.textContent = rowLabel;
            tr.appendChild(rowTh);

            this.labels.forEach((_, colIdx) => {
                const td = tr.insertCell();

                if (rowIdx === colIdx) {
                    td.className = "diagonal";
                    td.textContent = "—";
                } else {
                    const select = document.createElement("select");
                    OPTIONS.forEach(opt => {
                        const option = new Option(opt, opt);
                        select.add(option);
                    });
                    select.value = this.matrix[rowIdx][colIdx];
                    select.onchange = (e) => {
                        this.matrix[rowIdx][colIdx] = (e.target as HTMLSelectElement).value as Option;
                    };
                    td.appendChild(select);
                }
            });
        });
    }

    private async saveMatrix() {
        if (!this.currentUUID) {
            this.showStatus("ОШИБКА! Невозможно сохранить: отсутствует UUID", "red");
            return;
        }

        const intMatrix = this.matrix.map(row =>
            row.map(cell => REVERSE_MAP[cell])
        );

        try {
            await SaveMatrix(this.matrixFilename, intMatrix, this.currentUUID);
            this.showStatus("Сохранено!", "green");
        } catch (err: any) {
            this.showStatus("ОШИБКА! Невозможно сохранить", "red");
        }
    }

    private showStatus(msg: string, color: string) {
        this.statusEl.textContent = msg;
        this.statusEl.style.color = color;
        this.statusEl.style.fontSize = "200%";
        setTimeout(() => {
            if (this.statusEl.textContent === msg) this.statusEl.textContent = "";
        }, 4000);
    }
}

export function init() {
    new MatrixApp();
}