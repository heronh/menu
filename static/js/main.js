document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const logoutButton = document.getElementById('logoutButton'); // Will be in loaded navbar

    const tokenKey = 'authToken';
    const companyDataKey = 'companyData'; // To store company ID/details for easy access

    // --- Authentication ---
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('loginEmail').value;
            const password = document.getElementById('loginPassword').value;
            const errorDiv = document.getElementById('loginError');
            errorDiv.textContent = '';

            try {
                const response = await fetch('/api/v1/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password }),
                });
                const data = await response.json();
                if (response.ok) {
                    localStorage.setItem(tokenKey, data.token);
                    localStorage.setItem(companyDataKey, JSON.stringify(data.user)); // Store basic user/company info
                    window.location.href = '/web/company.html'; // Redirect to company admin page
                } else {
                    errorDiv.textContent = data.error || 'Login failed.';
                }
            } catch (error) {
                errorDiv.textContent = 'An error occurred during login.';
                console.error('Login error:', error);
            }
        });
    }

    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const name = document.getElementById('registerName').value;
            const email = document.getElementById('registerEmail').value;
            const password = document.getElementById('registerPassword').value;
            const passwordConfirmation = document.getElementById('registerPasswordConfirmation').value;
            const companyName = document.getElementById('registerCompanyName').value;
            const errorDiv = document.getElementById('registerError');
            errorDiv.textContent = '';

            if (password !== passwordConfirmation) {
                errorDiv.textContent = 'Passwords do not match.';
                return;
            }

            try {
                const response = await fetch('/api/v1/auth/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, email, password, password_confirmation: passwordConfirmation, company_name: companyName }),
                });
                const data = await response.json();
                if (response.ok) {
                    localStorage.setItem(tokenKey, data.token); // Assuming token is returned on successful registration
                    localStorage.setItem(companyDataKey, JSON.stringify(data.user));
                    window.location.href = '/web/company.html'; // Redirect
                } else {
                    errorDiv.textContent = data.error || 'Registration failed.';
                }
            } catch (error) {
                errorDiv.textContent = 'An error occurred during registration.';
                console.error('Registration error:', error);
            }
        });
    }

    function handleLogout() {
        localStorage.removeItem(tokenKey);
        localStorage.removeItem(companyDataKey);
        window.location.href = '/web/index.html'; // Redirect to home page
    }

    // --- Company Page Logic ---
    if (window.location.pathname.endsWith('company.html')) {
        const authToken = localStorage.getItem(tokenKey);
        if (!authToken) {
            window.location.href = '/web/index.html'; // Redirect if not logged in
            return; // Stop further execution on this page
        }

        loadTopShortcutBar();
        loadCompanyData(authToken);
        loadCompanyUsers(authToken);
        // Placeholder calls for other sections
        loadCompanyCategories(authToken);
        loadCompanyDishes(authToken);
        loadCompanyImages(authToken);

        const companyEditForm = document.getElementById('companyEditForm');
        if (companyEditForm) {
            companyEditForm.addEventListener('submit', handleUpdateCompany);
        }
    }


    async function loadTopShortcutBar() {
        const placeholder = document.getElementById('shortcutBarPlaceholder');
        if (placeholder) {
            try {
                const response = await fetch('/templates/layouts/top_shortcut_bar.html');
                if (response.ok) {
                    const html = await response.text();
                    placeholder.innerHTML = html;
                    // Re-attach logout button listener after loading navbar
                    const newLogoutButton = document.getElementById('logoutButton');
                    if (newLogoutButton) {
                        newLogoutButton.addEventListener('click', handleLogout);
                    }
                } else {
                    console.error('Failed to load top shortcut bar.');
                }
            } catch (error) {
                console.error('Error loading top shortcut bar:', error);
            }
        }
    }

    async function getAuthenticatedUserCompanyId() {
        const storedCompanyData = JSON.parse(localStorage.getItem(companyDataKey));
        if (storedCompanyData && storedCompanyData.company_id) {
            return storedCompanyData.company_id;
        }
        // Fallback or error if not found - this should ideally not happen if login/register is correct
        console.error("Company ID not found in local storage for current user.");
        return null;
    }


    async function loadCompanyData(token) {
        const companyId = await getAuthenticatedUserCompanyId();
        if (!companyId) {
            console.error("Cannot load company data without Company ID.");
            // Potentially redirect or show error to user
            return;
        }

        try {
            // Assuming /my endpoint correctly uses companyId from JWT claims on backend
            const response = await fetch(`/api/v1/companies/my`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const company = await response.json();
                document.getElementById('companyName').textContent = company.Name;
                document.getElementById('companyCNPJ').textContent = company.CNPJ || 'N/A';
                document.getElementById('companyZIPCode').textContent = company.ZIPCode || 'N/A';
                document.getElementById('companyStreet').textContent = company.Street || 'N/A';
                document.getElementById('companyNumber').textContent = company.Number || 'N/A';
                document.getElementById('companyNeighborhood').textContent = company.Neighborhood || 'N/A';
                document.getElementById('companyCity').textContent = company.City || 'N/A';
                document.getElementById('companyState').textContent = company.State || 'N/A';
                document.getElementById('companyLevel').textContent = company.Level;
                document.getElementById('companyActive').textContent = company.Active ? 'Yes' : 'No';

                // Populate edit form
                document.getElementById('editCompanyName').value = company.Name;
                document.getElementById('editCompanyCNPJ').value = company.CNPJ || '';
                document.getElementById('editCompanyZIPCode').value = company.ZIPCode || '';
                document.getElementById('editCompanyStreet').value = company.Street || '';
                document.getElementById('editCompanyNumber').value = company.Number || '';
                document.getElementById('editCompanyNeighborhood').value = company.Neighborhood || '';
                document.getElementById('editCompanyCity').value = company.City || '';
                document.getElementById('editCompanyState').value = company.State || '';
                document.getElementById('editCompanyLevel').value = company.Level;
                document.getElementById('editCompanyActive').value = company.Active ? "true" : "false";

                // Enable edit button
                const editCompanyButton = document.getElementById('editCompanyButton');
                if (editCompanyButton) editCompanyButton.disabled = false;

            } else if (response.status === 401) {
                handleLogout(); // Token expired or invalid
            }
            else {
                console.error('Failed to load company data:', response.status, await response.text());
                document.getElementById('companyName').textContent = 'Error loading data.';
            }
        } catch (error) {
            console.error('Error loading company data:', error);
            document.getElementById('companyName').textContent = 'Error loading data.';
        }
    }

    async function handleUpdateCompany(event) {
        event.preventDefault();
        const token = localStorage.getItem(tokenKey);
        const companyId = await getAuthenticatedUserCompanyId();
        if (!companyId || !token) {
            document.getElementById('companyEditError').textContent = 'Authentication error or company ID missing.';
            return;
        }

        const updatedData = {
            name: document.getElementById('editCompanyName').value,
            cnpj: document.getElementById('editCompanyCNPJ').value,
            zip_code: document.getElementById('editCompanyZIPCode').value,
            street: document.getElementById('editCompanyStreet').value,
            number: document.getElementById('editCompanyNumber').value,
            neighborhood: document.getElementById('editCompanyNeighborhood').value,
            city: document.getElementById('editCompanyCity').value,
            state: document.getElementById('editCompanyState').value,
            level: parseInt(document.getElementById('editCompanyLevel').value, 10),
            active: document.getElementById('editCompanyActive').value === "true",
        };

        // Filter out fields that were not changed or are empty, to send only partial updates.
        // GORM handles this well if pointers are used in request struct on backend and fields are omitted if nil.
        // For simplicity here, sending all, assuming backend handles empty strings vs. nulls appropriately.

        try {
            const response = await fetch(`/api/v1/companies/${companyId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(updatedData)
            });

            const result = await response.json();
            if (response.ok) {
                // Close modal
                const modalElement = document.getElementById('companyEditModal');
                const modalInstance = bootstrap.Modal.getInstance(modalElement);
                if (modalInstance) modalInstance.hide();

                // Refresh company data on page
                loadCompanyData(token);
                alert('Company data updated successfully!');
            } else {
                document.getElementById('companyEditError').textContent = result.error || 'Failed to update company data.';
            }
        } catch (error) {
            document.getElementById('companyEditError').textContent = 'An error occurred.';
            console.error('Error updating company:', error);
        }
    }


    async function loadCompanyUsers(token) {
        const companyId = await getAuthenticatedUserCompanyId();
        if (!companyId) return;

        const usersTableBody = document.getElementById('usersTableBody');
        if (!usersTableBody) return;

        try {
            const response = await fetch(`/api/v1/companies/${companyId}/users`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const users = await response.json();
                usersTableBody.innerHTML = ''; // Clear loading/previous
                if (users.length === 0) {
                    usersTableBody.innerHTML = '<tr><td colspan="4">No users found for this company.</td></tr>';
                    return;
                }
                users.forEach(user => {
                    const row = usersTableBody.insertRow();
                    row.innerHTML = `
                        <td>${user.name}</td>
                        <td>${user.email}</td>
                        <td>${user.privilege_name}</td>
                        <td>
                            <button class="btn btn-sm btn-info" onclick="alert('Edit user ${user.id} - Not implemented')">Edit</button>
                            <button class="btn btn-sm btn-danger" onclick="alert('Delete user ${user.id} - Not implemented')">Delete</button>
                        </td>
                    `;
                });
            } else if (response.status === 401) {
                handleLogout();
            } else {
                usersTableBody.innerHTML = '<tr><td colspan="4">Error loading users.</td></tr>';
                console.error('Failed to load company users:', response.status);
            }
        } catch (error) {
            usersTableBody.innerHTML = '<tr><td colspan="4">Error loading users.</td></tr>';
            console.error('Error loading company users:', error);
        }
    }

    // Placeholder functions for other sections
    function loadCompanyCategories(token) {
        const companyId = getAuthenticatedUserCompanyId();
        if (!companyId) return;
        const categoriesList = document.getElementById('categoriesList');
        if (categoriesList) categoriesList.innerHTML = '<p>Categories functionality not fully implemented.</p>';
        // Actual fetch: `/api/v1/companies/${companyId}/categories`
    }
    function loadCompanyDishes(token) {
        const companyId = getAuthenticatedUserCompanyId();
        if (!companyId) return;
        const dishesList = document.getElementById('dishesList');
        if (dishesList) dishesList.innerHTML = '<p>Dishes functionality not fully implemented.</p>';
        // Actual fetch: `/api/v1/companies/${companyId}/dishes`
    }
    function loadCompanyImages(token) {
        const companyId = getAuthenticatedUserCompanyId();
        if (!companyId) return;
        const imagesList = document.getElementById('imagesList');
        if (imagesList) imagesList.innerHTML = '<p>Images functionality not fully implemented.</p>';
        // Actual fetch: `/api/v1/companies/${companyId}/images`
    }

});
