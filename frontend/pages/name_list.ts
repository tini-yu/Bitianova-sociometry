import "../src/style.css";
import { ListAdd, ListGet, ListRemove, SaveTestDate } from "../wailsjs/go/main/App";

const DEBUG = true;

function log(...args: any[]) { if (DEBUG) console.log("[name_list]", ...args); }

function showStatus(statusEl: HTMLElement, message: string, type: "success" | "error" | "info" = "info") {
  statusEl.textContent = message;
  statusEl.className = `status ${type}`;
  setTimeout(() => {
    if (statusEl.textContent === message) {
      statusEl.textContent = "";
      statusEl.className = "status";
    }
  }, 2000);
}

async function refreshNames(listEl: HTMLElement, statusEl: HTMLElement) {
  try {
    log("refreshNames() start");
    const nameList: string[] = await ListGet();

    if (nameList.length === 0) {
      listEl.innerHTML = '<p class="empty">Список пуст.</p>';
      return;
    }

    listEl.innerHTML = nameList
      .map((name, index) => `
        <div class="word-item">
            <span class="word-text">${name}</span>
            <button class="delete-btn" data-index="${index}">×</button>
        </div>
      `).join("");

    listEl.querySelectorAll(".delete-btn").forEach((btn) => {
      btn.addEventListener("click", async (e) => {
        const target = e.target as HTMLElement;
        const index = parseInt(target.dataset.index || "-1");
        if (index === -1) return;

        try {
          await ListRemove(index);
          await refreshNames(listEl, statusEl);
          showStatus(statusEl, "Имя удалено", "success");
        } catch (err) {
          console.error("Delete error:", err);
          showStatus(statusEl, "Ошибка удаления", "error");
        }
      });
    });

    log("refreshNames() done, items:", nameList.length);
  } catch (err) {
    console.error("Failed to refresh names:", err);
    showStatus(statusEl, "Ошибка загрузки списка имён.", "error");
  }
}

export async function init() {
  try {
    log("init start");

    const inputEl = document.querySelector("#nameInput") as HTMLInputElement;
    const buttonEl = document.querySelector("#addButton") as HTMLButtonElement;
    const listEl = document.querySelector("#nameList") as HTMLDivElement;
    const statusEl = document.querySelector("#status") as HTMLParagraphElement;

    const dateInputEl = document.querySelector("#dateInput") as HTMLInputElement;
    const dateButtonEl = document.querySelector("#dateButton") as HTMLButtonElement;

    async function addDate() {
      const date = dateInputEl.value
      if (!date) {
        showStatus(statusEl, "Пожалуйста, введите дату.", "error");
        return;
      }
      try {
        await SaveTestDate(date)
        showStatus(statusEl, `Дата сохранена: ${date}`, "success");
      } catch (err) {
        showStatus(statusEl, `Неизвестная ошибка: ${err}`, "error");
      }
    }

    dateButtonEl.onclick = addDate;
    dateInputEl.onkeydown = (e) => {
      if (e.key === "Enter") addDate();
    };

    async function addName() {
      const name = inputEl.value.trim();
      if (!name) {
        showStatus(statusEl, "Пожалуйста, введите ФИО.", "error");
        return;
      }
      try {
        await ListAdd(name);
        inputEl.value = "";
        await refreshNames(listEl, statusEl);
        showStatus(statusEl, `Добавлен(а): ${name}`, "success");
      } catch (err) {
        showStatus(statusEl, "Имя уже в списке или ошибка.", "error");
      }
    }

    buttonEl.onclick = addName;
    inputEl.onkeydown = (e) => {
      if (e.key === "Enter") addName();
    };

    await refreshNames(listEl, statusEl);

    log("init done");
  } catch (err) {
    console.error("init error:", err);
  }
}