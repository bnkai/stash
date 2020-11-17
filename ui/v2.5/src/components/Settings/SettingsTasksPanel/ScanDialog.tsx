import React, { useState } from "react";
import { Button, Col, Form, Row } from "react-bootstrap";
import { useConfiguration } from "src/core/StashService";
import { Icon, Modal } from "src/components/Shared";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";

interface IScanDialogProps {
  onClose: (paths?: string[]) => void;
}

export const ScanDialog: React.FC<IScanDialogProps> = (
  props: IScanDialogProps
) => {
  const { data } = useConfiguration();

  const libraryPaths = data?.configuration.general.stashes.map((s) => s.path);

  const [paths, setPaths] = useState<string[]>([]);
  const [currentDirectory, setCurrentDirectory] = useState<string>("");

  function removePath(p: string) {
    setPaths(paths.filter((path) => path !== p));
  }

  function addPath(p: string) {
    if (p && !paths.includes(p)) {
      setPaths(paths.concat(p));
    }
  }

  return (
    <Modal
      show
      disabled={paths.length === 0}
      icon="pencil-alt"
      header="Select folders to scan"
      accept={{
        onClick: () => {
          props.onClose(paths);
        },
        text: "Scan",
      }}
      cancel={{
        onClick: () => props.onClose(),
        text: "Cancel",
        variant: "secondary",
      }}
    >
      <div className="dialog-container">
        {paths.map((p) => (
          <Row className="align-items-center mb-1">
            <Form.Label column xs={10}>
              {p}
            </Form.Label>
            <Col xs={2} className="d-flex justify-content-end">
              <Button
                className="ml-auto"
                size="sm"
                variant="danger"
                title="Delete"
                onClick={() => removePath(p)}
              >
                <Icon icon="minus" />
              </Button>
            </Col>
          </Row>
        ))}

        <FolderSelect
          currentDirectory={currentDirectory}
          setCurrentDirectory={(v) => setCurrentDirectory(v)}
          defaultDirectories={libraryPaths}
          appendButton={
            <Button
              variant="secondary"
              onClick={() => addPath(currentDirectory)}
            >
              <Icon icon="plus" />
            </Button>
          }
        />
      </div>
    </Modal>
  );
};
