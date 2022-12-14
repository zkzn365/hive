import React, { useState, FormEvent, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import type { FormDataType } from '@/common/interface';
import { useToast } from '@/hooks';
import { setNotice, getLoggedUserInfo } from '@/services';

const Index = () => {
  const toast = useToast();
  const { t } = useTranslation('translation', {
    keyPrefix: 'settings.notification',
  });
  const [formData, setFormData] = useState<FormDataType>({
    notice_switch: {
      value: false,
      isInvalid: false,
      errorMsg: '',
    },
  });

  const getProfile = () => {
    getLoggedUserInfo().then((res) => {
      setFormData({
        notice_switch: {
          value: res.notice_status === 1,
          isInvalid: false,
          errorMsg: '',
        },
      });
    });
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setNotice({
      notice_switch: formData.notice_switch.value,
    }).then(() => {
      toast.onShow({
        msg: t('update', { keyPrefix: 'toast' }),
        variant: 'success',
      });
    });
  };

  useEffect(() => {
    getProfile();
  }, []);
  return (
    <Form noValidate onSubmit={handleSubmit}>
      <Form.Group controlId="emailSend" className="mb-3">
        <Form.Label>{t('email.label')}</Form.Label>
        <Form.Check
          required
          type="checkbox"
          label={t('email.radio')}
          checked={formData.notice_switch.value}
          onChange={(e) => {
            setFormData({
              notice_switch: {
                value: e.target.checked,
                isInvalid: false,
                errorMsg: '',
              },
            });
          }}
        />
        <Form.Control.Feedback type="invalid">
          {formData.notice_switch.errorMsg}
        </Form.Control.Feedback>
      </Form.Group>

      <Button variant="primary" type="submit">
        {t('save', { keyPrefix: 'btns' })}
      </Button>
    </Form>
  );
};

export default React.memo(Index);
