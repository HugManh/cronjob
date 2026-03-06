const API_URL = '/api/v1/tasks';

document.addEventListener('DOMContentLoaded', () => {
    // Elements
    const taskGrid = document.getElementById('task-grid');
    const loadingSpinner = document.getElementById('loading');
    const emptyState = document.getElementById('empty-state');

    // Modal Elements
    const modalOverlay = document.getElementById('task-modal');
    const btnOpenModal = document.getElementById('open-new-task-modal');
    const btnCloseModal = document.getElementById('close-modal');
    const btnCancelModal = document.getElementById('cancel-modal');
    const taskForm = document.getElementById('task-form');
    const modalTitle = document.getElementById('modal-title');

    // Form Inputs
    const inputId = document.getElementById('task-id');
    const inputName = document.getElementById('task-name');
    const inputExecute = document.getElementById('task-execute');
    const inputMessage = document.getElementById('task-message');

    // Toast Element
    const toast = document.getElementById('toast');
    const toastMessage = document.getElementById('toast-message');

    // State
    let tasks = [];

    // Initialize
    fetchTasks();

    // Event Listeners
    btnOpenModal.addEventListener('click', () => openModal());
    btnCloseModal.addEventListener('click', closeModal);
    btnCancelModal.addEventListener('click', closeModal);
    taskForm.addEventListener('submit', handleTaskSubmit);

    // Close modal on outside click
    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) closeModal();
    });

    // --- Core Functions ---

    async function fetchTasks() {
        showLoading(true);
        try {
            const response = await fetch(`${API_URL}/`);
            if (!response.ok) throw new Error('Failed to fetch tasks');

            const result = await response.json();
            tasks = result.data || [];

            renderTasks();
        } catch (error) {
            showToast('Error loading tasks: ' + error.message, true);
        } finally {
            showLoading(false);
        }
    }

    function renderTasks() {
        taskGrid.innerHTML = '';

        if (tasks.length === 0) {
            taskGrid.style.display = 'none';
            emptyState.style.display = 'block';
            return;
        }

        taskGrid.style.display = 'grid';
        emptyState.style.display = 'none';

        tasks.forEach(task => {
            const isActive = task.active;
            const statusClass = isActive ? 'status-active' : 'status-inactive';
            const statusText = isActive ? 'Active' : 'Inactive';

            const card = document.createElement('div');
            card.className = 'task-card glass-panel';
            card.innerHTML = `
                <div class="task-header">
                    <h3 class="task-title">${escapeHTML(task.name)}</h3>
                    <span class="task-status ${statusClass}">${statusText}</span>
                </div>
                <div class="task-body">
                    <div class="cron-badge">${escapeHTML(task.execute)}</div>
                    <div class="task-message">${escapeHTML(task.message)}</div>
                </div>
                <div class="task-footer">
                    <span class="task-id" title="${task.id}">ID: ${task.id.substring(0, 8)}...</span>
                    <div class="task-actions">
                        <button class="btn btn-sm btn-secondary toggle-btn" data-id="${task.id}" data-active="${isActive}">
                            ${isActive ? 'Pause' : 'Resume'}
                        </button>
                        <button class="btn btn-sm btn-secondary edit-btn" data-id="${task.id}">Edit</button>
                        <button class="btn btn-sm btn-danger delete-btn" data-id="${task.id}">Delete</button>
                    </div>
                </div>
            `;

            taskGrid.appendChild(card);
        });

        // Attach event listeners to newly created buttons
        document.querySelectorAll('.edit-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const id = e.target.getAttribute('data-id');
                const task = tasks.find(t => t.id === id);
                if (task) openModal(task);
            });
        });

        document.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const id = e.target.getAttribute('data-id');
                if (confirm('Are you sure you want to delete this task?')) {
                    deleteTask(id);
                }
            });
        });

        document.querySelectorAll('.toggle-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const id = e.target.getAttribute('data-id');
                const currentState = e.target.getAttribute('data-active') === 'true';
                toggleTaskStatus(id, !currentState);
            });
        });
    }

    async function handleTaskSubmit(e) {
        e.preventDefault();

        const isEdit = !!inputId.value;
        const payload = {
            name: inputName.value.trim(),
            execute: inputExecute.value.trim(),
            message: inputMessage.value.trim(),
            active: "true"
        };

        try {
            let url = `${API_URL}/`;
            let method = 'POST';

            if (isEdit) {
                url = `${API_URL}/${inputId.value}`;
                method = 'PUT';
                // Active status persists usually via toggle but update endpoint accepts it
                const existingTask = tasks.find(t => t.id === inputId.value);
                if (existingTask) {
                    payload.active = existingTask.active ? "true" : "false";
                }
            }

            const response = await fetch(url, {
                method: method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!response.ok) throw new Error('Failed to save task');

            showToast(`Task ${isEdit ? 'updated' : 'created'} successfully!`);
            closeModal();
            fetchTasks();
        } catch (error) {
            showToast(error.message, true);
        }
    }

    async function deleteTask(id) {
        try {
            const response = await fetch(`${API_URL}/${id}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete task');

            showToast('Task deleted successfully!');
            fetchTasks();
        } catch (error) {
            showToast(error.message, true);
        }
    }

    async function toggleTaskStatus(id, newStatus) {
        try {
            const response = await fetch(`${API_URL}/${id}/active`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ active: newStatus })
            });

            if (!response.ok) throw new Error('Failed to change task status');

            fetchTasks(); // refresh to ensure sync
        } catch (error) {
            showToast(error.message, true);
        }
    }

    // --- UI Helpers ---

    function openModal(task = null) {
        if (task) {
            modalTitle.textContent = 'Edit Task';
            inputId.value = task.id;
            inputName.value = task.name;
            inputExecute.value = task.execute;
            inputMessage.value = task.message;
        } else {
            modalTitle.textContent = 'Create New Task';
            taskForm.reset();
            inputId.value = '';
        }

        modalOverlay.classList.add('active');
        inputName.focus();
    }

    function closeModal() {
        modalOverlay.classList.remove('active');
        taskForm.reset();
    }

    function showLoading(show) {
        loadingSpinner.style.display = show ? 'block' : 'none';
        if (show) {
            taskGrid.style.display = 'none';
            emptyState.style.display = 'none';
        }
    }

    let toastTimeout;
    function showToast(message, isError = false) {
        toastMessage.textContent = message;
        toast.style.backgroundColor = isError ? 'var(--danger-color)' : 'var(--success-bg)';
        toast.style.borderColor = isError ? 'rgba(248, 81, 73, 0.4)' : 'rgba(46, 160, 67, 0.3)';

        toast.classList.add('show');

        clearTimeout(toastTimeout);
        toastTimeout = setTimeout(() => {
            toast.classList.remove('show');
        }, 3000);
    }

    function escapeHTML(str) {
        if (!str) return '';
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
});
