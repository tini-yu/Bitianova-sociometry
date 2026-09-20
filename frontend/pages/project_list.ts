import "./project_list.css";
import { GetProjectList, ChangeProjectDir, CreateProject, DeleteProject, IsExistProjectDir, GetExpectedProjectDir, ConvertNameToProjectDir } from "../wailsjs/go/main/App";
import { setProjectLoaded } from "../src/menu_state"

class ProjectApp {
    private listEl!: HTMLDivElement;
    private selectedElement: HTMLElement | null = null;

    constructor() {
        this.waitForDOMAndInit();
    }

    private async waitForDOMAndInit() {
        this.listEl = document.querySelector("#projectList") as HTMLDivElement;
        this.bindButtons();
        await this.refreshList();
    }

    private bindButtons() {
        document.getElementById("loadButton")!.addEventListener("click", () => this.handleLoadClick());
        document.getElementById("createButton")!.addEventListener("click", () => this.handleCreateClick());
        document.getElementById("deleteButton")!.addEventListener("click", () => this.handleDeleteClick());
    }

    private async refreshList() {
        try {
            const projects: string[] = await GetProjectList();

            if (projects.length === 0) {
                this.listEl.innerHTML = '<p class="empty">Список проектов пуст.</p>';
                this.selectedElement = null;
                return;
            }

            let html = "";
            for (const name of projects) {
                html += `
                    <div class="project-name" data-project="${name}">
                        <span class="name-text">${name}</span>
                    </div>`;
            }
            this.listEl.innerHTML = html;

            this.addClickHandlers();

        } catch (err) {
            console.error("Failed to refresh names:", err);
            this.listEl.innerHTML = `<p style="color:red">Ошибка загрузки списка</p>`;
        }
    }

    private addClickHandlers() {
        this.listEl.querySelectorAll(".project-name").forEach(item => {
            const el = item as HTMLElement;

            el.onclick = () => this.selectItem(el);
        });
    }

    private selectItem(item: HTMLElement) {
        // Снимаем выделение со всех
        this.listEl.querySelectorAll(".project-name").forEach(i =>
            i.classList.remove("selected")
        );

        item.classList.add("selected");
        this.selectedElement = item;
    }

    private handleLoadClick() {
        if (!this.selectedElement) {
            this.showStatus("Необходимо выбрать проект", "red");
            return;
        }

        const projectName = this.selectedElement.dataset.project;

        if (!projectName) {
            this.showStatus("Ошибка: имя проекта не найдено", "red");
            return;
        }

        this.loadProject(projectName);
    }

    private async loadProject(projectName: string) {
        try {
            const exists = await IsExistProjectDir(projectName);

            if (!exists) {
                this.showStatus(`Директория проекта "${projectName}" не найдена`, "red");
                setProjectLoaded(false);   // ← глобально
                return;
            }

            const pname: string = await ChangeProjectDir(projectName);
            console.log(pname)

            this.showStatus(`Проект "${projectName}" успешно загружен`, "green");
            this.checkCurrentProject();
            // setProjectLoaded(true);        // ← включаем меню

        } catch (err: any) {
            console.error(err);
            this.showStatus(`Ошибка загрузки: ${err.message || err}`, "red");
            setProjectLoaded(false);
        }
    }

    private async checkCurrentProject() {
        try {
            const current = await GetExpectedProjectDir();
            if (current) {
                const exists = await IsExistProjectDir(current);
                setProjectLoaded(exists);
            }
        } catch (_) { }
    }

    // СОЗДАНИЕ НОВОГО ПРОЕКТА:
    private async handleCreateClick() {

        const modal = document.getElementById("createProjectModal") as HTMLDivElement;
        const input = document.getElementById("newProjectName") as HTMLInputElement;

        input.value = "";
        modal.style.display = "flex";

        // Фокус на поле ввода
        setTimeout(() => input.focus(), 100);

        // Enter = создать
        input.onkeydown = (e) => {
            if (e.key === "Enter") this.confirmCreate();
        };

        document.getElementById("confirmCreateBtn")!.onclick = () => this.confirmCreate();
        document.getElementById("cancelCreateBtn")!.onclick = () => this.closeModal();
        ;
    }

    private async confirmCreate() {
        const input = document.getElementById("newProjectName") as HTMLInputElement;
        let name = input.value.trim();

        if (!name) {
            alert("Имя проекта не может быть пустым!");
            input.focus();
            return;
        }

        // === Проверка на валидность имени папки в Windows ===
        const validation = this.validateProjectName(name);
        if (!validation.isValid) {
            alert(validation.errorMessage);
            input.focus();
            return;
        }

        try {
            // this.showStatus(`Создание проекта "${name}"...`, "blue");

            await CreateProject(name);

            this.closeModal();
            this.showStatus(`Проект "${name}" успешно создан`, "green");
            await this.loadProject(name);
            await this.refreshList();

        } catch (err: any) {
            console.error(err);
            this.showStatus(`Ошибка создания проекта: ${err.message || err}`, "red");
        }
    }

    // ==================== ВАЛИДАЦИЯ ИМЕНИ ====================
    private validateProjectName(name: string): { isValid: boolean; errorMessage: string } {

        // 1. Проверка на запрещённые символы
        const invalidChars = /[<>:"\/\\|?*]/g;
        if (invalidChars.test(name)) {
            return {
                isValid: false,
                errorMessage: "Имя содержит недопустимые символы:\n< > : \" / \\ | ? *"
            };
        }

        // 2. Не может заканчиваться точкой или пробелом
        if (name.endsWith('.') || name.endsWith(' ')) {
            return {
                isValid: false,
                errorMessage: "Имя не может заканчиваться точкой или пробелом"
            };
        }

        // 3. Зарезервированные имена Windows
        const reservedNames = [
            "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
            "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5",
            "LPT6", "LPT7", "LPT8", "LPT9"
        ];

        const upperName = name.toUpperCase();
        if (reservedNames.includes(upperName)) {
            return {
                isValid: false,
                errorMessage: `Имя "${name}" является зарезервированным системой Windows`
            };
        }

        // 4. Максимальная длина (оставляем запас)
        if (name.length > 200) {
            return {
                isValid: false,
                errorMessage: "Имя проекта слишком длинное (максимум 200 символов)"
            };
        }

        return { isValid: true, errorMessage: "" };
    }

    private closeModal() {
        const modal = document.getElementById("createProjectModal") as HTMLDivElement;
        modal.style.display = "none";
    }

    // ====================== DELETE ======================
    private handleDeleteClick() {
        if (!this.selectedElement) {
            this.showStatus("Необходимо выбрать проект", "red");
            return;
        }

        const projectName = this.selectedElement.dataset.project!;
        this.showDeleteModal(projectName);
    }

    private async showDeleteModal(projectName: string) {
        const confirmBtn = document.getElementById("confirmDeleteBtn") as HTMLButtonElement;
        const cancelBtn = document.getElementById("cancelDeleteBtn") as HTMLButtonElement;

        // === СБРОС КНОПОК В ИСХОДНОЕ СОСТОЯНИЕ ===
        confirmBtn.style.display = "inline-block";
        confirmBtn.textContent = "Удалить";
        cancelBtn.textContent = "Отмена";
        // Сбрасываем обработчики (на всякий случай)
        confirmBtn.onclick = null;
        cancelBtn.onclick = null;

        const currentProject = await GetExpectedProjectDir();

        // === ЗАПРЕТ УДАЛЕНИЯ ТЕКУЩЕГО ПРОЕКТА ===
        if (await ConvertNameToProjectDir(projectName) === currentProject) {

            const modal = document.getElementById("deleteProjectModal") as HTMLDivElement;
            const textEl = document.getElementById("deleteConfirmText") as HTMLParagraphElement;

            textEl.innerHTML = `
                    <span style="color: #c70000; font-weight: bold;">
                        Невозможно удалить текущий проект!<br>
                    </span>
                    Сначала загрузите другой проект.
                `;
            modal.style.display = "flex";

            // Кнопки скрываем или делаем только "Отмена"
            document.getElementById("confirmDeleteBtn")!.style.display = "none";
            document.getElementById("cancelDeleteBtn")!.textContent = "Понятно";
            document.getElementById("cancelDeleteBtn")!.onclick = () => this.closeDeleteModal();
            return;
        }

        const modal = document.getElementById("deleteProjectModal") as HTMLDivElement;
        const textEl = document.getElementById("deleteConfirmText") as HTMLParagraphElement;

        textEl.textContent = `Вы действительно хотите удалить проект "${projectName}"?`;
        modal.style.display = "flex";

        document.getElementById("confirmDeleteBtn")!.onclick = () => this.confirmDelete(projectName);
        document.getElementById("cancelDeleteBtn")!.onclick = () => this.closeDeleteModal();
    }

    private async confirmDelete(projectName: string) {
        try {
            // this.showStatus(`Удаление проекта "${projectName}"...`, "blue");
            this.closeDeleteModal();

            await DeleteProject(projectName);

            this.showStatus(`Проект "${projectName}" успешно удалён`, "green");

            // После удаления текущего проекта — блокируем меню
            const exists = await IsExistProjectDir(projectName);
            setProjectLoaded(exists);           // обычно false

            this.selectedElement = null;

            this.refreshList();

        } catch (err: any) {
            let errorMsg = "Неизвестная ошибка";

            if (err && typeof err === 'string') errorMsg = err;
            else if (err?.message) errorMsg = err.message;
            else if (err?.toString) errorMsg = err.toString();

            this.showStatus(`Ошибка удаления: ${errorMsg}`, "red");
            setProjectLoaded(false);
        }
    }

    private closeDeleteModal() {
        const modal = document.getElementById("deleteProjectModal") as HTMLDivElement;
        modal.style.display = "none";
    }

    private showStatus(msg: string, color: string) {
        const statusEl = document.getElementById("status") as HTMLDivElement;
        if (!statusEl) return;

        statusEl.textContent = msg;
        statusEl.style.color = color;

        setTimeout(() => {
            if (statusEl.textContent === msg) statusEl.textContent = "";
        }, 5000);
    }
}

export function init() {
    new ProjectApp();
}