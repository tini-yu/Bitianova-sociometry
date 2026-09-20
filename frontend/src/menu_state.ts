// src/menuState.ts
let isProjectLoaded = false;

export function setProjectLoaded(loaded: boolean): void {
    isProjectLoaded = loaded;
    updateSidebarState();
}

export function getProjectLoaded(): boolean {
    return isProjectLoaded;
}

// Обновляет состояние бокового меню
function updateSidebarState(): void {
    const sidebarLinks = document.querySelectorAll('.sidebar a');
    
    sidebarLinks.forEach(link => {
        const page = (link as HTMLElement).dataset.page;
        
        // Эти страницы всегда доступны
        if (page === 'project_list') {
            link.classList.remove('disabled');
            return;
        }

        if (isProjectLoaded) {
            link.classList.remove('disabled');
        } else {
            link.classList.add('disabled');
        }
    });
}

// Вызывается при инициализации
export function initMenuState(): void {
    updateSidebarState();
}