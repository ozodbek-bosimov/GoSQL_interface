// Global variables
let currentEditId = null;
let searchTimeout = null;

// Tab switching
function showTab(tabName, event) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab').forEach(tab => {
        tab.classList.remove('active');
    });

    // Show selected tab
    document.getElementById(tabName + '-tab').classList.add('active');
    if (event && event.target) {
        event.target.classList.add('active');
    } else {
        // Fallback: find the tab button by text content
        document.querySelectorAll('.tab').forEach(tab => {
            if (tab.textContent.toLowerCase().includes(tabName.toLowerCase())) {
                tab.classList.add('active');
            }
        });
    }

    // Load schema if schema tab is selected
    if (tabName === 'schema') {
        fetchSchema();
    }
}

// Fetch records from API
async function fetchRecords(search = '') {
    try {
        const url = `${apiEndpoint}?search=${encodeURIComponent(search)}`;
        const response = await fetch(url);
        const data = await response.json();

        renderTable(data);
        document.getElementById('total-count').textContent = data.length;
    } catch (error) {
        console.error('Error fetching records:', error);
        showMessage('Failed to load data', 'error');
    }
}

// Render table rows
function renderTable(data) {
    const tbody = document.getElementById('table-body');
    tbody.innerHTML = '';

    if (data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="100" style="text-align: center; padding: 20px;">No data found</td></tr>';
        return;
    }

    data.forEach(record => {
        const row = document.createElement('tr');
        row.innerHTML = generateTableRow(record);
        tbody.appendChild(row);
    });
}

// Generate table row HTML based on entity type
function generateTableRow(record) {
    if (entityType === 'users') {
        return `
            <td>${record.id}</td>
            <td>${record.name || '-'}</td>
            <td>${record.login}</td>
            <td>${record.role}</td>
            <td>${new Date(record.created_at).toLocaleDateString('uz-UZ')}</td>
            <td class="action-buttons">
                <button class="btn btn-sm btn-primary" onclick="openEditModal(${record.id})">Edit</button>
                <button class="btn btn-sm btn-danger" onclick="deleteRecord(${record.id})">Delete</button>
            </td>
        `;
    } else if (entityType === 'interfaces') {
        return `
            <td>${record.id}</td>
            <td>${record.name}</td>
            <td>${record.ip}</td>
            <td>${record.mac}</td>
            <td>${record.mtu}</td>
            <td><span class="badge ${record.status ? 'badge-success' : 'badge-danger'}">${record.status ? 'Active' : 'Inactive'}</span></td>
            <td><span class="badge ${record.ip_type ? 'badge-success' : 'badge-danger'}">${record.ip_type ? 'Static' : 'Dynamic'}</span></td>
            <td class="action-buttons">
                <button class="btn btn-sm btn-primary" onclick="openEditModal(${record.id})">Edit</button>
                <button class="btn btn-sm btn-danger" onclick="deleteRecord(${record.id})">Delete</button>
            </td>
        `;
    } else if (entityType === 'acl') {
        return `
            <td>${record.id}</td>
            <td>${record.src_ip}</td>
            <td>${record.dst_ip}</td>
            <td>${record.protocol.toUpperCase()}</td>
            <td>${record.src_port || '-'}</td>
            <td>${record.dst_port || '-'}</td>
            <td><span class="badge ${record.action === 'permit' || record.action === 'permit+reflect' ? 'badge-success' : 'badge-danger'}">${record.action}</span></td>
            <td class="action-buttons">
                <button class="btn btn-sm btn-primary" onclick="openEditModal(${record.id})">Edit</button>
                <button class="btn btn-sm btn-danger" onclick="deleteRecord(${record.id})">Delete</button>
            </td>
        `;
    }
}

// Search with debouncing
function searchRecords() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        const searchTerm = document.getElementById('search').value;
        fetchRecords(searchTerm);
    }, 300);
}

// Open modal for adding new record
function openAddModal() {
    currentEditId = null;
    document.getElementById('modal-title').textContent = 'Add New';
    generateFormFields();
    document.getElementById('modal').style.display = 'flex';
}

// Open modal for editing record
async function openEditModal(id) {
    currentEditId = id;
    document.getElementById('modal-title').textContent = 'Edit';

    try {
        const response = await fetch(`${apiEndpoint}/${id}`);
        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to fetch record');
        }
        const record = await response.json();
        generateFormFields(record);
        document.getElementById('modal').style.display = 'flex';
    } catch (error) {
        console.error('Error fetching record:', error);
        showMessage('Failed to load record', 'error');
    }
}

// Generate form fields
function generateFormFields(data = {}) {
    const formFields = document.getElementById('form-fields');
    formFields.innerHTML = '';

    fields.forEach(field => {
        const formGroup = document.createElement('div');
        formGroup.className = 'form-group';

        const label = document.createElement('label');
        label.textContent = field.label;
        label.htmlFor = field.name;
        formGroup.appendChild(label);

        let input;
        if (field.type === 'select') {
            input = document.createElement('select');
            input.id = field.name;
            input.name = field.name;
            input.required = field.required;

            field.options.forEach(option => {
                const opt = document.createElement('option');
                opt.value = option.value;
                opt.textContent = option.label;
                input.appendChild(opt);
            });

            if (data[field.name]) {
                input.value = data[field.name];
            }
        } else if (field.type === 'checkbox') {
            const wrapper = document.createElement('div');
            wrapper.className = 'checkbox-wrapper';

            input = document.createElement('input');
            input.type = 'checkbox';
            input.id = field.name;
            input.name = field.name;
            input.checked = data[field.name] || false;

            const checkboxLabel = document.createElement('label');
            checkboxLabel.htmlFor = field.name;
            checkboxLabel.textContent = 'Yes';

            wrapper.appendChild(input);
            wrapper.appendChild(checkboxLabel);
            formGroup.appendChild(wrapper);
            formFields.appendChild(formGroup);
            return;
        } else {
            input = document.createElement('input');
            input.type = field.type;
            input.id = field.name;
            input.name = field.name;
            input.required = field.required;

            if (field.placeholder) {
                input.placeholder = field.placeholder;
            }

            if (data[field.name] !== undefined && data[field.name] !== null) {
                input.value = data[field.name];
            }

            // For password field in edit mode, make it optional
            if (field.type === 'password' && currentEditId) {
                input.required = false;
            }
        }

        formGroup.appendChild(input);
        formFields.appendChild(formGroup);
    });
}

// Close modal
function closeModal() {
    document.getElementById('modal').style.display = 'none';
    currentEditId = null;
}

// Save record (create or update)
async function saveRecord() {
    const formData = {};

    fields.forEach(field => {
        const element = document.getElementById(field.name);
        if (element.type === 'checkbox') {
            formData[field.name] = element.checked;
        } else if (element.type === 'number') {
            formData[field.name] = parseInt(element.value);
        } else if (element.value) {
            formData[field.name] = element.value;
        } else if (!field.required) {
            formData[field.name] = null;
        }
    });

    try {
        let url = apiEndpoint;
        let method = 'POST';

        if (currentEditId) {
            url = `${apiEndpoint}/${currentEditId}`;
            method = 'PUT';
        }

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to save record');
        }

        closeModal();
        fetchRecords();
        showMessage(currentEditId ? 'Successfully updated' : 'Successfully created', 'success');
    } catch (error) {
        console.error('Error saving record:', error);
        showMessage(error.message || 'Failed to save', 'error');
    }
}

// Delete record
async function deleteRecord(id) {
    if (!confirm('Are you sure you want to delete this record?')) {
        return;
    }

    try {
        const response = await fetch(`${apiEndpoint}/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to delete record');
        }

        fetchRecords();
        showMessage('Successfully deleted', 'success');
    } catch (error) {
        console.error('Error deleting record:', error);
        showMessage('Failed to delete', 'error');
    }
}

// Fetch and display schema
async function fetchSchema() {
    console.log('Fetching schema from:', `${apiEndpoint}/schema`);
    try {
        const response = await fetch(`${apiEndpoint}/schema`);
        console.log('Schema response status:', response.status);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const schema = await response.json();
        console.log('Schema data received:', schema);

        const schemaBody = document.getElementById('schema-body');
        if (!schemaBody) {
            console.error('Element #schema-body not found!');
            return;
        }

        schemaBody.innerHTML = '';

        if (!schema || schema.length === 0) {
            schemaBody.innerHTML = '<tr><td colspan="4" style="text-align: center; padding: 20px;">No schema data available</td></tr>';
            return;
        }

        schema.forEach(column => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td><strong>${column.name}</strong></td>
                <td>${column.type}</td>
                <td>${column.nullable ? 'Yes' : 'No'}</td>
                <td>${column.primary ? 'Yes' : 'No'}</td>
            `;
            schemaBody.appendChild(row);
        });

        console.log('Schema table populated with', schema.length, 'columns');
    } catch (error) {
        console.error('Error fetching schema:', error);
        const schemaBody = document.getElementById('schema-body');
        if (schemaBody) {
            schemaBody.innerHTML = `<tr><td colspan="4" style="text-align: center; padding: 20px; color: #e74c3c;">Error loading schema: ${error.message}</td></tr>`;
        }
    }
}

// Show success/error message
function showMessage(message, type) {
    // Create message element
    const msgDiv = document.createElement('div');
    msgDiv.className = type === 'success' ? 'success-message' : 'error-message';
    msgDiv.textContent = message;
    msgDiv.style.position = 'fixed';
    msgDiv.style.top = '20px';
    msgDiv.style.right = '20px';
    msgDiv.style.zIndex = '9999';
    msgDiv.style.padding = '15px 20px';
    msgDiv.style.borderRadius = '4px';
    msgDiv.style.boxShadow = '0 2px 5px rgba(0,0,0,0.2)';

    document.body.appendChild(msgDiv);

    // Remove after 3 seconds
    setTimeout(() => {
        msgDiv.remove();
    }, 3000);
}

// Close modal on outside click
window.onclick = function(event) {
    const modal = document.getElementById('modal');
    if (event.target === modal) {
        closeModal();
    }
}
