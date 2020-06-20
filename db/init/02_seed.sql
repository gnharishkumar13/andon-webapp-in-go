INSERT INTO public.workcenters ( wc_name, current_product, wc_status, escalation_level, status_set_at ) VALUES 
('Assembly Line 1', 'Widgets', 2, 1, NOW() - interval '5 minutes'),
('Assembly Line 2', 'Widgets', 1, 1, NOW()),
('Assembly Line 3', 'Widgets', 0, 0, NOW()),
('Roll Furnace 1', 'Widgets', 0, 0, NOW()),
('Roll Furnace 2', 'Widgets', 0, 0, NOW()),
('Roll Furnace 3', 'Widgets', 0, 0, NOW()),
('Roll Furnace 4', 'Widgets', 0, 0, NOW()),
('Roll Furnace 5', 'Widgets', 0, 0, NOW()),
('Transfer Press 1', 'Widgets', 0, 0, NOW()),
('Transfer Press 2', 'Widgets', 0, 0, NOW()),
('Transfer Press 3', 'Widgets', 0, 0, NOW()),
('Progressive Press 1', 'Widgets', 0, 0, NOW()),
('Progressive Press 2', 'Widgets', 0, 0, NOW()),
('Progressive Press 3', 'Widgets', 0, 0, NOW()),
('Progressive Press 4', 'Widgets', 0, 0, NOW());

--- password: admin
INSERT INTO public.users ( username, password )
VALUES ('admin', 'S84yvdRCuzUlqG1SQWs3vHNaaSB5FbG8RXPAzMIti6VW+v4sGRYQzsVDyoOvkRQYbsGGxXQqkichoxZvXadoEA==');

--- password: nonadmin
INSERT INTO public.users ( username, password )
VALUES ('nonadmin', 'hO4zvS+e5qQXU4VIzVrmuFdpBGfqyz22pwy9JmgFNOvyckHNRY6V/C0erUCP2SdVR31Y6ZiJy8YiqEMG9+N51A==');


INSERT INTO public.roles (role) VALUES
	('admin');

INSERT INTO public.users_roles (user_id, role_id)
	SELECT users.id as user_id, roles.id as role_id
	FROM users, roles
	WHERE users.username = 'admin'
		AND roles.role = 'admin';
