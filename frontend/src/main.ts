import "./style.css";

interface PageLink extends HTMLElement {
  dataset: {
    page: string;
  };
}

// находим все файлы страниц автоматически
// as: 'raw' заставляет Vite загружать HTML как строку (текст), а не как код
const htmlModules = import.meta.glob('../pages/*.html', { as: 'raw' });
// загружаем скрипты как динамические модули
const scriptModules = import.meta.glob('../pages/*.ts');

async function loadPage(pageName: string): Promise<void> {
  try {
    const container = document.getElementById("page-container");
    if (!container) throw new Error("Container not found");

    // путь до html страниц
    const htmlPath = `../pages/${pageName}.html`;
    
    if (!htmlModules[htmlPath]) {
        throw new Error(`HTML file not found: ${htmlPath}`);
    }

    // функция импорта для получения контента
    const htmlContent = await htmlModules[htmlPath](); 
    container.innerHTML = htmlContent;

    const scriptPath = `../pages/${pageName}.ts`;

    if (scriptModules[scriptPath]) {
        const module: any = await scriptModules[scriptPath]();
        //init вызывается при каждом переходе на страницу
        if (module.init && typeof module.init === 'function') {
            module.init();
        }
    }

    // менюшка
    document.querySelectorAll(".sidebar a").forEach((link) => {
      link.classList.toggle("active", (link as PageLink).dataset.page === pageName);
    });

  } catch (err) {
    const container = document.getElementById("page-container");
    if (container) {
      container.innerHTML = `<p style='color: red;'>Ошибка: ${(err as Error).message}</p>`;
    }
    console.error("Load page error:", err);
  }
}

function setupMenuHandlers(): void {
  document.querySelectorAll(".sidebar a").forEach((link) => {
    link.addEventListener("click", (e: Event) => {
      e.preventDefault();
      const page = (link as PageLink).dataset.page;
      if(page) loadPage(page);
    });
  });
}

document.addEventListener("DOMContentLoaded", () => {
  setupMenuHandlers();
  loadPage("name_list");
});