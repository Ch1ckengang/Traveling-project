import { useEffect, useState, useCallback } from 'react';
import { adminGetUsers, adminUpdateUserStatus, adminUpdateUserRole } from '../../services/adminService';
import './AdminPage.css';

const AdminUsersPage = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const [search, setSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState('');

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const res = await adminGetUsers({ page, limit: 10, role: roleFilter, search });
      setUsers(res.data || []);
      setMeta(res.meta || {});
    } catch {
      setUsers([]);
    }
    setLoading(false);
  }, [page, search, roleFilter]);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  const handleToggleStatus = async (userId, currentActive) => {
    const action = currentActive ? 'khóa' : 'kích hoạt';
    if (!window.confirm(`Bạn có chắc muốn ${action} tài khoản này?`)) return;
    try {
      await adminUpdateUserStatus(userId, !currentActive);
      fetchUsers();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const handleChangeRole = async (userId, newRole) => {
    if (!window.confirm(`Thay đổi role thành "${newRole}"?`)) return;
    try {
      await adminUpdateUserRole(userId, newRole);
      fetchUsers();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const roleLabels = { customer: 'Khách hàng', staff: 'Nhân viên', admin: 'Quản trị' };

  const formatDate = (d) => {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('vi-VN');
  };

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Quản lý Users</h1>
      </div>

      <div className="admin-filters">
        <input
          type="text"
          placeholder="Tìm tên, email, SĐT..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="filter-input"
        />
        <select value={roleFilter} onChange={(e) => { setRoleFilter(e.target.value); setPage(1); }} className="filter-select">
          <option value="">Tất cả role</option>
          <option value="customer">Khách hàng</option>
          <option value="staff">Nhân viên</option>
          <option value="admin">Quản trị</option>
        </select>
      </div>

      {loading ? (
        <div className="admin-loading">Đang tải...</div>
      ) : (
        <>
          <div className="admin-table-wrapper">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Tên</th>
                  <th>Email</th>
                  <th>SĐT</th>
                  <th>Role</th>
                  <th>Trạng thái</th>
                  <th>Ngày tạo</th>
                  <th>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {users.length === 0 ? (
                  <tr><td colSpan="8" className="empty-row">Không có user nào</td></tr>
                ) : users.map(user => (
                  <tr key={user.id}>
                    <td>{user.id}</td>
                    <td>{user.name}</td>
                    <td>{user.email}</td>
                    <td>{user.phone || '—'}</td>
                    <td>
                      <select
                        value={user.role}
                        onChange={(e) => handleChangeRole(user.id, e.target.value)}
                        className="role-select"
                      >
                        <option value="customer">Khách hàng</option>
                        <option value="staff">Nhân viên</option>
                        <option value="admin">Quản trị</option>
                      </select>
                    </td>
                    <td>
                      <span className={`status-dot ${user.is_active ? 'active' : 'inactive'}`}>
                        {user.is_active ? 'Hoạt động' : 'Bị khóa'}
                      </span>
                    </td>
                    <td>{formatDate(user.created_at)}</td>
                    <td className="actions">
                      <button
                        className={`btn-sm ${user.is_active ? 'btn-cancel-action' : 'btn-confirm'}`}
                        onClick={() => handleToggleStatus(user.id, user.is_active)}
                      >
                        {user.is_active ? 'Khóa' : 'Mở khóa'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {meta.total_pages > 1 && (
            <div className="admin-pagination">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Trước</button>
              <span>Trang {page} / {meta.total_pages} ({meta.total} users)</span>
              <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Sau →</button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default AdminUsersPage;
