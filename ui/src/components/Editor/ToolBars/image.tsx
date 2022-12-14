import { FC, useEffect, useState, memo } from 'react';
import { Button, Form, Modal, Tab, Tabs } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import { Modal as AnswerModal } from '@/components';
import ToolItem from '../toolItem';
import { IEditorContext } from '../types';
import { uploadImage } from '@/services';

const Image: FC<IEditorContext> = ({ editor }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'editor' });

  const loadingText = `![${t('image.uploading')}...]()`;

  const item = {
    label: 'image',
    keyMap: ['Ctrl-G'],
    tip: `${t('image.text')} (Ctrl+G)`,
  };
  const [currentTab, setCurrentTab] = useState('localImage');
  const [visible, setVisible] = useState(false);
  const [link, setLink] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
    type: '',
  });

  const [imageName, setImageName] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
  });
  const verifyImageSize = (files: FileList) => {
    if (files.length === 0) {
      return false;
    }
    const filteredFiles = Array.from(files).filter(
      (file) => file.type.indexOf('image') === -1,
    );

    if (filteredFiles.length > 0) {
      AnswerModal.confirm({
        content: t('image.only_image'),
      });
      return false;
    }
    const filteredImages = Array.from(files).filter(
      (file) => file.size / 1024 / 1024 > 4,
    );

    if (filteredImages.length > 0) {
      AnswerModal.confirm({
        content: t('image.max_size'),
      });
      return false;
    }
    return true;
  };
  const upload = (
    files: FileList,
  ): Promise<{ url: string; name: string }[]> => {
    const promises = Array.from(files).map(async (file) => {
      const url = await uploadImage({ file, type: 'post' });

      return {
        name: file.name,
        url,
      };
    });

    return Promise.all(promises);
  };
  function dragenter(_, e) {
    e.stopPropagation();
    e.preventDefault();
  }

  function dragover(_, e) {
    e.stopPropagation();
    e.preventDefault();
  }
  const drop = async (_, e) => {
    if (!editor) {
      return;
    }
    const fileList = e.dataTransfer.files;

    const bool = verifyImageSize(fileList);

    if (!bool) {
      return;
    }

    const startPos = editor.getCursor();
    const endPos = { ...startPos, ch: startPos.ch + loadingText.length };

    editor.replaceSelection(loadingText);
    const urls = await upload(fileList);

    const text = urls.map(({ name, url }) => {
      return `![${name}](${url})`;
    });

    editor.replaceRange(text.join('\n'), startPos, endPos);
  };

  const paste = async (_, event) => {
    const clipboard = event.clipboardData;

    const bool = verifyImageSize(clipboard.files);

    if (bool) {
      event.preventDefault();
      editor.setOption('readOnly', true);
      const startPos = editor.getCursor('');
      const endPos = { ...startPos, ch: startPos.ch + loadingText.length };

      editor.replaceSelection(loadingText);
      const urls = await upload(clipboard.files);
      const text = urls.map(({ name, url }) => {
        return `![${name}](${url})`;
      });

      editor.replaceRange(text.join('\n'), startPos, endPos);

      editor.setOption('readOnly', false);
      return;
    }

    const htmlStr = clipboard.getData('text/html');
    const imgRegex = /<img([\s\S]*?) src\s*=\s*(['"])([\s\S]*?)\2([^>]*)>/;

    if (!htmlStr.match(imgRegex)) {
      return;
    }
    event.preventDefault();

    const newHtml = new DOMParser()
      .parseFromString(
        htmlStr.replace(
          /<img([\s\S]*?) src\s*=\s*(['"])([\s\S]*?)\2([^>]*)>/gi,
          `<p>\n![${t('image.text')}]($3)\n</p>`,
        ),
        'text/html',
      )
      .querySelector('body')?.innerText as string;

    editor.replaceSelection(newHtml);
  };
  const handleClick = () => {
    if (!editor) {
      return;
    }
    if (!link.value) {
      setLink({ ...link, isInvalid: true });
      return;
    }
    setLink({ ...link, type: '' });

    const text = `![${imageName.value}](${link.value})`;

    editor.replaceSelection(text);

    setVisible(false);

    editor.focus();
    setLink({ ...link, value: '' });
    setImageName({ ...imageName, value: '' });
  };
  useEffect(() => {
    editor.on('dragenter', dragenter);
    editor.on('dragover', dragover);
    editor.on('drop', drop);
    editor.on('paste', paste);
    return () => {
      editor.off('dragenter', dragenter);
      editor.off('dragover', dragover);
      editor.off('drop', drop);
      editor.off('paste', paste);
    };
  }, [editor]);

  useEffect(() => {
    if (link.value && link.type === 'drop') {
      handleClick();
    }
  }, [link.value]);

  const addLink = () => {
    if (!editor) {
      return;
    }
    const text = editor?.getSelection();

    setImageName({ ...imageName, value: text });

    setVisible(true);
  };

  const onUpload = async (e) => {
    if (!editor) {
      return;
    }
    const files = e.target?.files || [];
    const bool = verifyImageSize(files);

    if (!bool) {
      return;
    }

    uploadImage({ file: e.target.files[0], type: 'post' }).then((url) => {
      setLink({ ...link, value: url });
    });
  };

  const onHide = () => setVisible(false);
  const onExited = () => editor?.focus();

  const handleSelect = (tab) => {
    setCurrentTab(tab);
  };
  return (
    <ToolItem {...item} onClick={addLink}>
      <Modal
        show={visible}
        onHide={onHide}
        onExited={onExited}
        fullscreen="sm-down">
        <Modal.Header closeButton>
          <h5 className="mb-0">{t('image.add_image')}</h5>
        </Modal.Header>
        <Modal.Body>
          <Tabs onSelect={handleSelect}>
            <Tab eventKey="localImage" title={t('image.tab_image')}>
              <Form className="mt-3" onSubmit={handleClick}>
                <Form.Group controlId="editor.imgLink" className="mb-3">
                  <Form.Label>
                    {t('image.form_image.fields.file.label')}
                  </Form.Label>
                  <Form.Control
                    type="file"
                    onChange={onUpload}
                    isInvalid={currentTab === 'localImage' && link.isInvalid}
                  />

                  <Form.Control.Feedback type="invalid">
                    {t('image.form_image.fields.file.msg.empty')}
                  </Form.Control.Feedback>
                </Form.Group>

                <Form.Group controlId="editor.imgDescription" className="mb-3">
                  <Form.Label>
                    {t('image.form_image.fields.description.label')}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={imageName.value}
                    onChange={(e) =>
                      setImageName({ ...imageName, value: e.target.value })
                    }
                    isInvalid={imageName.isInvalid}
                  />
                </Form.Group>
              </Form>
            </Tab>
            <Tab eventKey="remoteImage" title={t('image.tab_url')}>
              <Form className="mt-3" onSubmit={handleClick}>
                <Form.Group controlId="editor.imgUrl" className="mb-3">
                  <Form.Label>
                    {t('image.form_url.fields.url.label')}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={link.value}
                    onChange={(e) =>
                      setLink({ ...link, value: e.target.value })
                    }
                    isInvalid={currentTab === 'remoteImage' && link.isInvalid}
                  />
                  <Form.Control.Feedback type="invalid">
                    {t('image.form_url.fields.url.msg.empty')}
                  </Form.Control.Feedback>
                </Form.Group>

                <Form.Group controlId="editor.imgName" className="mb-3">
                  <Form.Label>
                    {t('image.form_url.fields.name.label')}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={imageName.value}
                    onChange={(e) =>
                      setImageName({ ...imageName, value: e.target.value })
                    }
                    isInvalid={imageName.isInvalid}
                  />
                </Form.Group>
              </Form>
            </Tab>
          </Tabs>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="link" onClick={() => setVisible(false)}>
            {t('image.btn_cancel')}
          </Button>
          <Button variant="primary" onClick={handleClick}>
            {t('image.btn_confirm')}
          </Button>
        </Modal.Footer>
      </Modal>
    </ToolItem>
  );
};

export default memo(Image);
