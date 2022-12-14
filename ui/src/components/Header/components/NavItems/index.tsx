import { FC, memo } from 'react';
import { Nav, Dropdown } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { Link, NavLink } from 'react-router-dom';

import type * as Type from '@/common/interface';
import { Avatar, Icon } from '@/components';

interface Props {
  redDot: Type.NotificationStatus | undefined;
  userInfo;
  logOut: () => void;
}

const Index: FC<Props> = ({ redDot, userInfo, logOut }) => {
  const { t } = useTranslation();
  return (
    <>
      <Nav.Link
        as={NavLink}
        to="/users/notifications/inbox"
        className="icon-link d-flex align-items-center justify-content-center p-0 me-3 position-relative">
        <div className="text-white text-opacity-75">
          <Icon name="bell-fill" className="fs-4" />
        </div>
        {(redDot?.inbox || 0) > 0 && <div className="unread-dot bg-danger" />}
      </Nav.Link>

      <Nav.Link
        as={Link}
        to="/users/notifications/achievement"
        className="icon-link d-flex align-items-center justify-content-center p-0 me-3 position-relative">
        <div className="text-white text-opacity-75">
          <Icon name="trophy-fill" className="fs-4" />
        </div>
        {(redDot?.achievement || 0) > 0 && (
          <div className="unread-dot bg-danger" />
        )}
      </Nav.Link>

      <Dropdown align="end">
        <Dropdown.Toggle
          variant="success"
          id="dropdown-basic"
          as="a"
          className="no-toggle pointer">
          <Avatar size="36px" avatar={userInfo?.avatar} searchStr="s=96" />
        </Dropdown.Toggle>

        <Dropdown.Menu>
          <Dropdown.Item href={`/users/${userInfo.username}`}>
            {t('header.nav.profile')}
          </Dropdown.Item>
          <Dropdown.Item href="/users/settings/profile">
            {t('header.nav.setting')}
          </Dropdown.Item>
          {userInfo?.is_admin ? (
            <Dropdown.Item href="/admin">{t('header.nav.admin')}</Dropdown.Item>
          ) : null}
          {redDot?.can_revision ? (
            <Dropdown.Item href="/review" className="position-relative">
              {t('header.nav.review')}
              {redDot?.revision > 0 && (
                <span className="position-absolute top-50 translate-middle-y end-0 me-3 p-2 bg-danger border border-light rounded-circle">
                  <span className="visually-hidden">New Review</span>
                </span>
              )}
            </Dropdown.Item>
          ) : null}
          <Dropdown.Divider />
          <Dropdown.Item onClick={logOut}>
            {t('header.nav.logout')}
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    </>
  );
};

export default memo(Index);
