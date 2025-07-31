import { Route } from "wouter";
import { Toaster } from "react-hot-toast";
import { SidebarProvider } from "./components/ui/sidebar";
import SidebarNav from "./components/layouts/sidebar-nav";
import { AuthProvider } from "./contexts/AuthContext";
import AuthTest from "./pages/AuthTest";
import Login from "./modules/auth/Login";
import Register from "./modules/auth/Register";
import ProtectedRoute from "./modules/auth/ProtectedRoute";
import Files from "./modules/content/files";
import MediaFiles from "./modules/content/media-files";
import UserSettings from "./modules/settings/UserSettings";
import AdminRoute from "./modules/auth/AdminRoute";
import UserManagementPage from "./modules/master-data/user-management/user-management-page";
import DashboardPage from "./modules/dashboard/pages/dashboard-page";

function AppContent() {
  return (
    <>
      <Toaster
        toastOptions={{
          duration: 5000,
        }}
      />
      <Route path="/login">{<Login />}</Route>
      <Route path="/register">{<Register />}</Route>
      <Route path="/auth-test">{<AuthTest />}</Route>

      <ProtectedRoute>
        <SidebarProvider>
          <div className="flex min-h-screen w-screen">
            <SidebarNav />
            <main className="m-4 h-auto w-full">
              <Route path="/">{<DashboardPage />}</Route>
              <Route path="/images">{<MediaFiles mediaType="image" />}</Route>
              <Route path="/documents">{<MediaFiles mediaType="document" />}</Route>
              <Route path="/videos">{<MediaFiles mediaType="video" />}</Route>
              <Route path="/audio">{<MediaFiles mediaType="audio" />}</Route>
              <Route path="/settings">{<UserSettings />}</Route>
              <Route path="/admin/user-management">
                <AdminRoute>
                  <UserManagementPage />
                </AdminRoute>
              </Route>
            </main>
          </div>
        </SidebarProvider>
      </ProtectedRoute>
    </>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
