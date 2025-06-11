document.addEventListener('DOMContentLoaded', () => {
    // 获取所有需要的 DOM 元素
    const productForm = document.getElementById('productForm');
    const stockInForm = document.getElementById('stockInForm');
    const stockOutForm = document.getElementById('stockOutForm'); // 出库表单
    const searchKeywordInput = document.getElementById('searchKeyword');
    const inventoryTable = document.getElementById('inventoryTable');
    const recordTable = document.getElementById('recordTable');
    const inProductSelect = document.getElementById('inProductSelect');
    const outProductSelect = document.getElementById('outProductSelect'); // 出库商品选择

    // 统计卡片元素
    const totalProductsEl = document.getElementById('totalProducts');
    const todayInEl = document.getElementById('todayIn');
    const todayOutEl = document.getElementById('todayOut');
    const lowStockEl = document.getElementById('lowStock');

    // --- 数据获取与更新函数 ---

    async function fetchProducts() {
        try {
            const res = await fetch('/products');
            if (!res.ok) throw new Error('获取商品失败');
            const data = await res.json();
            updateProductSelects(data);
            updateInventoryTable(data);
        } catch (error) {
            console.error('商品加载失败:', error);
            showToast('商品加载失败', 'danger');
        }
    }

    async function fetchRecords() {
        try {
            const res = await fetch('/records');
            if (!res.ok) throw new Error('获取记录失败');
            const data = await res.json();
            updateRecordTable(data);
        } catch (error) {
            console.error('记录加载失败:', error);
            showToast('记录加载失败', 'danger');
        }
    }

    async function fetchStats() {
        try {
            const res = await fetch('/stats');
            if (!res.ok) throw new Error('获取统计失败');
            const stats = await res.json();
            updateStatsDashboard(stats);
        } catch (error) {
            console.error('统计数据加载失败:', error);
            showToast('统计数据加载失败', 'danger');
        }
    }

    // --- UI 更新函数 ---

    function updateStatsDashboard(stats) {
        totalProductsEl.textContent = stats.totalProducts || 0;
        todayInEl.textContent = stats.todayIn || 0;
        todayOutEl.textContent = stats.todayOut || 0;
        lowStockEl.textContent = stats.lowStock || 0;
    }

    function updateProductSelects(products) {
        if (!products) return;
        [inProductSelect, outProductSelect].forEach(select => {
            select.innerHTML = '<option value="">请选择商品</option>';
            products.forEach(p => {
                const option = document.createElement('option');
                option.value = p.ID;
                option.textContent = `${p.Name} (ID: ${p.ID})`;
                select.appendChild(option);
            });
        });
    }

    function updateInventoryTable(products) {
       if (!products) return;
       inventoryTable.innerHTML = '';
       products.forEach(p => {
           const stockStatus = p.Stock <= p.SafetyStock ? 'status-low' : 'status-normal';
           const statusText = p.Stock <= p.SafetyStock ? '低库存' : '正常';
           inventoryTable.innerHTML += `
               <tr>
                   <td>${p.ID}</td>
                   <td>${p.Name}</td>
                   <td>${p.Stock}</td>
                   <td>${p.SafetyStock}</td>
                   <td><span class="status-indicator ${stockStatus}">${statusText}</span></td>
               </tr>
           `;
       });
    }

    function updateRecordTable(records) {
        if (!records) return;
        recordTable.innerHTML = '';
        records.slice().reverse().forEach(r => {
            recordTable.innerHTML += `
                <tr>
                    <td>${r.ProductID}</td>
                    <td>${r.Operation === 'in' ? '入库' : '出库'}</td>
                    <td>${r.Quantity}</td>
                    <td>${r.Operator}</td>
                    <td>${new Date(r.Timestamp * 1000).toLocaleString()}</td>
                </tr>
            `;
        });
    }


    searchKeywordInput.addEventListener('input', async () => {
        const keyword = searchKeywordInput.value.trim();
        try {
            // 根据是否有关键字决定请求的URL
            const url = keyword ? `/inventory?keyword=${keyword}` : '/products';
            const res = await fetch(url);
            if (!res.ok) throw new Error('搜索失败');
            const data = await res.json();
            updateInventoryTable(data);
        } catch (error) {
            console.error('搜索失败:', error);
        }
    });

    // 添加新商品
    productForm.addEventListener('submit', async e => {
        e.preventDefault();
        const btn = productForm.querySelector('button');
        btn.disabled = true;
        const product = {
            Name: document.getElementById('productName').value,
            SafetyStock: +document.getElementById('safetyStock').value
        };
        try {
            const res = await fetch('/products', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(product)
            });
            if (res.ok) {
                await fetchProducts();
                await fetchStats();
                productForm.reset();
                showToast('商品添加成功', 'success');
            } else {
                const errorData = await res.text();
                showToast(`商品添加失败: ${errorData}`, 'danger');
            }
        } catch (err) {
            showToast('网络错误', 'danger');
        } finally {
            btn.disabled = false;
        }
    });

    // 入库
    stockInForm.addEventListener('submit', async e => {
        e.preventDefault();
        const btn = stockInForm.querySelector('button');
        btn.disabled = true;
        const productID = inProductSelect.value;
        const quantity = +document.getElementById('inQuantity').value;
        const operator = document.getElementById('inOperator').value;
        try {
            const res = await fetch(`/stock/in?product_id=${productID}&quantity=${quantity}&operator=${operator}`, {
                method: 'POST'
            });
            if (res.ok) {
                await fetchProducts();
                await fetchRecords();
                await fetchStats();
                stockInForm.reset();
                showToast('入库成功', 'success');
            } else {
                const errorData = await res.text();
                showToast(`入库失败: ${errorData}`, 'danger');
            }
        } catch (err) {
            showToast('网络错误', 'danger');
        } finally {
            btn.disabled = false;
        }
    });


    stockOutForm.addEventListener('submit', async e => {
        e.preventDefault();
        const btn = stockOutForm.querySelector('button');
        btn.disabled = true;
        const productID = outProductSelect.value;
        const quantity = +document.getElementById('outQuantity').value;
        const operator = document.getElementById('outOperator').value;
        try {
            const res = await fetch(`/stock/out?product_id=${productID}&quantity=${quantity}&operator=${operator}`, {
                method: 'POST'
            });
            if (res.ok) {
                await fetchProducts();
                await fetchRecords();
                await fetchStats();
                stockOutForm.reset();
                showToast('出库成功', 'success');
            } else {
                const errorData = await res.text();
                showToast(`出库失败: ${errorData}`, 'danger');
            }
        } catch (err) {
            showToast('网络错误', 'danger');
        } finally {
            btn.disabled = false;
        }
    });

    // 弹窗提示
    function showToast(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = `toast align-items-center text-white bg-${type} border-0 position-fixed bottom-0 end-0 m-3`;
        toast.setAttribute('role', 'alert');
        toast.setAttribute('aria-live', 'assertive');
        toast.setAttribute('aria-atomic', 'true');
        toast.innerHTML = `
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
        `;
        document.body.appendChild(toast);
        const bsToast = new bootstrap.Toast(toast, { delay: 3000 });
        bsToast.show();
        toast.addEventListener('hidden.bs.toast', () => toast.remove());
    }

    // --- 初始加载和定时刷新 ---
    function fetchAllData() {
        fetchProducts();
        fetchRecords();
        fetchStats();
    }

    fetchAllData(); // 页面初次加载
    setInterval(fetchAllData, 30000); // 每30秒刷新
});